package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Music struct {
	MusicID    string           `json:"musichash"`
	Name       string           `json:"name"`
	Artist     string           `json:"artist"`
	Length     string           `json:"length"`
	Points     []CopyrightStock `json:"points"`
	Authorized []Certi          `json:"authorized"`
}
type CopyrightStock struct {
	Code  string `json:"key"`
	Owner string `json:"owner"`
	Sale  bool   `json:"sale"`
}
type Certi struct {
	TokenHash string `json:"key"`
	MusicID   string `json:"musicID"`
	Owner     string `json:"owner"`
	Contract  string `json:"contract"`
	Expired   bool   `json:"expired"`
}
type Artist struct {
	UserId string           `json:"key"`
	Name   string           `json:"name"`
	Points []CopyrightStock `json:"points"`
	Musics []Music          `json:"musics"`
}
type User struct {
	UserID string           `json:"key"`
	Name   string           `json:"name"`
	Points []CopyrightStock `json:"points"`
}
type SecondCreator struct {
	UserID        string  `json:"key"`
	Name          string  `json:"name"`
	Certification []Certi `json:"certifications"`
}

// 객체 = 구조체 + 메서드
type MusicAsset struct {
}

// 객체를 shim에 전달할때 필요한 함수 호출
// t *MusicAsset : MusicAsset에 속하는 메서드
// 입력값 형식 : APIstub
// 출력값 형식 : peer.Resoponse
func (t *MusicAsset) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *MusicAsset) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	// CLI나 Application에서 전달되는 Argument 가져오기
	fn, arg := APIstub.GetFunctionAndParameters()
	// Invoke : set & get > 원장에 데이터를 저장 / 조회 (key:value 형식)
	var result string
	var err error

	if fn == "register" {
		result, err = t.registerMusic(APIstub, arg)
	} else if fn == "querystock" {
		result, err = t.querySingleStock(APIstub, arg)
	} else if fn == "queryallstock" {
		result, err = t.queryAllStock(APIstub, arg)
	} else if fn == "sell" {
		result, err = t.SellMusic(APIstub, arg)
	} else if fn == "buy" {
		result, err = t.buyMusic(APIstub, arg)
	} else {
		return shim.Error("Not supporte chaincode function!!")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

// 음원 창작자 ID 생성
func (t *MusicAsset) registerAritist(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("invalid number of Arguments!")
	}
	id := args[0]
	name := args[1]
	artist := Artist{UserId: id, Name: name, Points: nil, Musics: nil}

	artistAsBytes, _ := json.Marshal(artist)
	err := APIstub.PutState(id, artistAsBytes)

	if err != nil {
		return "", fmt.Errorf("failed to register new artist: %s", err)
	}

	return "", err
}

// 2차 창작자 ID 생성
func (t *MusicAsset) register2ndCreator(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("invalid number of Arguments!")
	}
	id := args[0]
	name := args[1]
	_2ndCreator := SecondCreator{UserID: id, Name: name, Certification: nil}

	_2ndAsBytes, _ := json.Marshal(_2ndCreator)
	err := APIstub.PutState(id, _2ndAsBytes)

	if err != nil {
		return "", fmt.Errorf("failed to register secondary creator: %s", err)
	}

	return "", err
}

// 단일 지분 토큰 조회
func (t *MusicAsset) querySingleStock(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}

	pointAsBytes, _ := APIstub.GetState(args[0])
	return string(pointAsBytes), nil
}

// 음원에 대한 전체 지분 토큰 조회
func (t *MusicAsset) queryAllStock(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	// musicID로 조회 & 해당 음악의 지분 리스트 추출
	musicAsBytes, _ := APIstub.GetState(args[0])
	target := Music{}
	err := json.Unmarshal(musicAsBytes, &target)
	list := target.Points

	startKey := list[0].Code
	endKey := list[len(list)-1].Code

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"코드\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"상태\":")

		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

	return buffer.String(), nil
}

// 제목, 아티스트, ID, 음악길이, 지분 생성 개수를 입력받아 지분 등록
func (t *MusicAsset) registerMusic(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 5 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}
	// 원작자가 요청한 개수만큼 지분 생성
	value, _ := strconv.Atoi(args[4])
	head := args[0] + args[2]

	points := t.IssueStock(APIstub, head, value, args[2])

	var music = Music{Name: args[0], Artist: args[1], MusicID: args[2], Length: args[3], Points: points, Authorized: nil}
	musicAsBytes, _ := json.Marshal(music)
	err := APIstub.PutState(args[2], musicAsBytes)

	if err != nil {
		return "", fmt.Errorf("failed to set asset: %s", err)
	}

	return string(musicAsBytes), nil
}

// 토큰 코드번호 생성 (value만큼 point 발행)
func (t *MusicAsset) IssueStock(APIstub shim.ChaincodeStubInterface, head string, value int, owner string) []CopyrightStock {
	points := make([]CopyrightStock, value)
	for i := 1; i <= value; i++ {
		code := head + int2code(i)
		points[i-1] = CopyrightStock{Code: code, Owner: owner, Sale: true}
		pointAsBytes, _ := json.Marshal(code)
		_ = APIstub.PutState(code, pointAsBytes)
	}
	return points
}

// 코드값 생성 함수
func int2code(n int) string {
	code := strconv.Itoa(n)
	if n < 10 {
		return "0000" + code
	} else if n < 100 {
		return "000" + code
	} else if n < 1000 {
		return "00" + code
	} else if n < 10000 {
		return "0" + code
	} else {
		return code
	}
}

// 나의 음원 리스트
func (t *MusicAsset) showMusicList(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("invalid number of Arguments!")
	}
	id := args[0]
	artistAsBytes, _ := APIstub.GetState(id)
	artist := Artist{}
	err := json.Unmarshal(artistAsBytes, &artist)
	name := artist.Name

	fmt.Println("%s님의 등록 곡 리스트\n", name)

	list := artist.Musics
	// 원작자의 Musics 리스트에 1개 이상의 음원이 등록되었다면
	if list != nil {
		startKey := list[0].MusicID
		endKey := list[len(list)-1].MusicID

		resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
		if err != nil {
			return "", err
		}
		defer resultsIterator.Close()

		var buffer bytes.Buffer

		bArrayMemberAlreadyWritten := false
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return "", err
			}

			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"음원 ID\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"내용\":")

			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}

		fmt.Printf("- 음원 리스트:\n%s\n", buffer.String())
	}
	return "", err
}

// 나의 음원 이용자 조회하기
func (t *MusicAsset) queryCerti(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("invalid number of Arguments!")
	}
	musicId := args[0]
	musicAsBytes, _ := APIstub.GetState(musicId)
	music := Music{}
	err := json.Unmarshal(musicAsBytes, &music)

	list := music.Authorized
	// 원작자의 Musics 리스트에 1개 이상의 음원이 등록되었다면
	if list != nil {
		startKey := list[0].TokenHash
		endKey := list[len(list)-1].TokenHash

		resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
		if err != nil {
			return "", err
		}
		defer resultsIterator.Close()

		var buffer bytes.Buffer

		bArrayMemberAlreadyWritten := false
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return "", err
			}

			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"음원 Hash\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"상태\":")

			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}

		fmt.Printf("- 음원 리스트:\n%s\n", buffer.String())
	} else {
		fmt.Println("허가권 발행 이력이 없습니다.\n")
	}
	return "", err
}

// 테스트 파일 생성
func (t *MusicAsset) initLedger(APIstub shim.ChaincodeStubInterface) peer.Response {
	musics := []Music{
		Music{MusicID: "BP8211", Name: "HYLT", Artist: "BlackPink", Length: "2:56",
			Points: []CopyrightStock{
				{Code: "BP82111", Owner: "Alice", Sale: false},
				{Code: "BP82112", Owner: "Alice", Sale: false},
				{Code: "BP82113", Owner: "Alice", Sale: false}},
			Authorized: nil},
		Music{MusicID: "NAVIS7", Name: "Savage", Artist: "Aespa", Length: "3:48",
			Points: []CopyrightStock{
				{Code: "NAVIS71", Owner: "Alice", Sale: false},
				{Code: "NAVIS72", Owner: "Bob", Sale: true},
				{Code: "NAVIS73", Owner: "Bob", Sale: false}},
			Authorized: nil},
		Music{MusicID: "BTSSTB", Name: "Permission to Dance", Artist: "BTS", Length: "3:23",
			Points: []CopyrightStock{
				{Code: "BTSSTB1", Owner: "Alice", Sale: false},
				{Code: "BTSSTB2", Owner: "Bob", Sale: true},
				{Code: "BTSSTB3", Owner: "Charlie", Sale: true}},
			Authorized: nil}}

	i := 0

	for i < len(musics) {
		fmt.Println("i is ", i)
		musicAsBytes, _ := json.Marshal(musics[i])
		APIstub.PutState("CAR"+strconv.Itoa(i), musicAsBytes)
		fmt.Println("Added", musics[i])
		i = i + 1
	}

	return shim.Success(nil)
}

// 지분 판매 등록
func (t *MusicAsset) sellMusic(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("invalid number of Arguments!")
	}

	pointAsBytes, err := APIstub.GetState(args[0]) //저작권 코드로 조회
	copyright := CopyrightStock{}

	json.Unmarshal(pointAsBytes, &copyright)
	copyright.Sale = true // 판매중 상태를 true로 변경

	pointAsBytes, _ = json.Marshal(copyright)
	APIstub.PutState(args[0], pointAsBytes)

	return "", err
}

// 지분 구매
func (t *MusicAsset) buyMusic(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("invalid number of Arguments!")
	}

	pointAsBytes, err := APIstub.GetState(args[0]) //저작권 코드로 조회
	copyright := CopyrightStock{}

	json.Unmarshal(pointAsBytes, &copyright)
	copyright.Owner = args[1] // 소유주 변경
	copyright.Sale = false    // 판매중 상태를 true로 변경

	pointAsBytes, _ = json.Marshal(copyright)
	APIstub.PutState(args[0], pointAsBytes)

	return "", err
}

// 허가권 구매 : onwerId, musicId, myId, duration
func (t *MusicAsset) buyCertification(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect number of arguments !!")
	}

	// 토큰 해시 생성
	value := TokenHash("certi")

}

func main() {
	err := shim.Start(new(MusicAsset)) // shim에 chaincode 객체를 인자로 전달
	if err != nil {
		fmt.Printf("error creating new Smart Contract: %s", err)
	}
}
