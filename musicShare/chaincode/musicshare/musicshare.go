package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	peer "github.com/hyperledger/fabric/protos/peer"
)

type Music struct {
	MusicID     string     `json:"musichash"`
	Name        string     `json:"name"`
	Artist      string     `json:"artist"`
	Length      string     `json:"length"`
	EntireStake int        `json:"entire_stake"`
	Contracts   []Contract `json:"contracts"`
}
type Contract struct {
	MusicID string `json:"musicID"` // Q : 또 필요한지
	Owner   string `json:"owner"`   // set
	Buyer   string `json:"buyer"`
	Stake   int    `json:"stake"`   // set
	Date    string `json:"date"`    // set
	Expired string `json:"expired"` // set
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
	} else if fn == "set" {
		result, err = t.setContract(APIstub, arg)
	} else if fn == "fill" {
		result, err = t.fillInContract(APIstub, arg)
	} else if fn == "query" {
		result, err = t.queryContract(APIstub, arg)
	} else if fn == "expire" {
		result, err = t.expire(APIstub, arg)
	} else if fn == "share" {
		result, err = t.shareProfit(APIstub, arg)
	} else {
		return shim.Error("Not supporte chaincode function!!")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func (t *MusicAsset) initMusic(APIstub shim.ChaincodeStubInterface) (string, error) {
	musics := []Music{
		Music{
			MusicID: "0001", Name: "Savage", Artist: "Aespa", Length: "3:59", EntireStake: 200, Contracts: nil},
		Music{
			MusicID: "0002", Name: "작은것들을위한시", Artist: "BTS", Length: "3:12", EntireStake: 300, Contracts: nil},
		Music{
			MusicID: "0003", Name: "술한잔해요", Artist: "지아", Length: "4:13", EntireStake: 50, Contracts: nil},
		Music{
			MusicID: "0004", Name: "Tempo", Artist: "EXO", Length: "3:46", EntireStake: 150, Contracts: nil},
		Music{
			MusicID: "0005", Name: "Enemy", Artist: "Imagine Dragons", Length: "2:53", EntireStake: 200, Contracts: nil}}

	for i := 0; i < len(musics); i++ {
		fmt.Println("i is ", i)
		musicAsBytes, _ := json.Marshal(musics[i])
		APIstub.PutState(musics[i].MusicID, musicAsBytes)
		fmt.Println("Added", musics[i])
	}
	return "", nil
}

// 음원 등록: musicID, 제목, 아티스트, 음악길이, 지분 생성 개수를 입력받음
func (t *MusicAsset) registerMusic(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 5 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}
	// 원작자가 요청한 개수만큼 지분 생성
	stake, _ := strconv.Atoi(args[3])

	var music = Music{MusicID: args[0], Name: args[1], Artist: args[2], Length: args[3], EntireStake: stake}
	musicAsBytes, _ := json.Marshal(music)
	err := APIstub.PutState(args[0], musicAsBytes)

	if err != nil {
		return "", fmt.Errorf("failed to set asset: %s", err)
	}

	return string(musicAsBytes), nil
}

// 계약서 생성 : musicId, myId, stake, duration
// date 추가하기
func (t *MusicAsset) setContract(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect number of arguments !!")
	}
	// 음원명에 해당하는 Music 갖고오기 ?
	musicAsBytes, err := APIstub.GetState(args[0])
	music := Music{}
	json.Unmarshal(musicAsBytes, &music)
	stake, _ := strconv.Atoi(args[2])
	owner := music.Artist

	// 만료기한 설정

	now := time.Now()
	convDays, _ := strconv.Atoi(args[3])
	expired := now.AddDate(0, 0, convDays).Format("2006-01-02")
	contract := Contract{MusicID: args[0], Owner: owner, Buyer: args[1], Stake: stake, Date: "", Expired: expired}

	music.Contracts = append(music.Contracts, contract)
	musicAsBytes, _ = json.Marshal(music)

	err = APIstub.PutState(args[0], musicAsBytes)
	return "the contract has successfully set", err
}

// 이용자 계약서 내용 채우기  : 음원ID -> 계약내용 추가
// 후원자 음원 후원         - stake N
// 2차 창작자 음원 대여 구매 - stake 0
// musicID, Buyer, Stake, Date(자동생성)
func (t *MusicAsset) fillInContract(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}

	// musicID로 조회
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}
	stake2int, _ := strconv.Atoi(args[2])
	entire2int := music.EntireStake

	err := json.Unmarshal(musicAsBytes, &music)

	// list := music.Contracts
	// preset := list[0]

	// 후원자는 entireStake 줄이기
	if stake2int != 0 {
		// 남은 stake 체크
		if music.EntireStake-entire2int < 0 {
			return "", fmt.Errorf("No more stake. Out of EntireStake.")
		}

		music.EntireStake -= entire2int
	}

	// 계약내용 채우기
	date := time.Now().Format("2006-01-02")

	contract := Contract{MusicID: args[0], Owner: args[1], Buyer: args[2], Stake: stake2int, Date: date, Expired: args[4]}

	music.Contracts = append(music.Contracts, contract)

	err = APIstub.PutState(args[0], musicAsBytes)
	return "", err
}

// 계약서 조회 : musicId
func (t *MusicAsset) queryContract(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	// musicID로 조회 & 해당 음악의 지분 리스트 추출
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}
	err := json.Unmarshal(musicAsBytes, &music)

	list := music.Contracts
	if list == nil {
		return "the list is empty !!", nil
	}
	list_d := music.Contracts // 후원자 list
	list_c := music.Contracts // 2차창작자 list

	// 후원자 : 2차창작자 분리
	for _, l := range list {
		if l.Stake != 0 {
			list_d = append(list_d, l)
		} else {
			list_c = append(list_c, l)
		}
	}

	var buffer bytes.Buffer

	if list_d != nil {
		buffer.WriteString(music.MusicID)
		buffer.WriteString(" 후원 리스트\n")
		for _, d := range list_d {
			buffer.WriteString("소유주: ")
			buffer.WriteString(d.Owner)
			buffer.WriteString(" 지분: ")
			stake2str := strconv.Itoa(d.Stake)
			buffer.WriteString(stake2str)
			buffer.WriteString(" 만료일: ")
			buffer.WriteString(d.Expired)
			buffer.WriteString("\n")
		}
	}
	if list_c != nil {
		buffer.WriteString(music.MusicID)
		buffer.WriteString(" 저작권 대여 리스트\n")
		for _, c := range list_c {
			buffer.WriteString("소유주: ")
			buffer.WriteString(c.Owner)
			buffer.WriteString(" 만료일: ")
			buffer.WriteString(c.Expired)
			buffer.WriteString("\n")
		}
	}

	fmt.Printf("- 계약내용조회:\n%s\n", buffer.String())

	return buffer.String(), err
}

// 수익 분배 : musicId, profit
// 원작자 수익 -> 지분에 따라 후원자에게 분배
// 2차 창작자 수익 -> 계약 내용에 따라 원작자에게 분배
func (t *MusicAsset) shareProfit(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	// musicID로 조회 & 해당 음악의 지분 리스트 추출
	musicAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		return "Invalid music ID !!", err
	}
	music := Music{}
	whole, _ := strconv.Atoi(args[1]) // 총수익
	err = json.Unmarshal(musicAsBytes, &music)

	list := music.Contracts
	if list == nil {
		return "the list is empty !!", nil
	}
	list_d := []Contract{} // 유효 후원자 list

	wholeStake := music.EntireStake
	now := time.Now().Format("2006-01-02")
	now2int, _ := strconv.Atoi(now)

	// 지분
	for _, l := range list {
		expired2int, _ := strconv.Atoi(l.Expired)
		left := expired2int - now2int
		if l.Stake != 0 && (left >= 0) {
			list_d = append(list_d, l)
			wholeStake += l.Stake
		}
	}

	var buffer bytes.Buffer

	if list_d != nil {
		buffer.WriteString(music.MusicID)
		buffer.WriteString(" 수익분배표\n")
		for _, d := range list_d {
			share := whole * d.Stake / wholeStake
			buffer.WriteString("소유주: ")
			buffer.WriteString(d.Owner)
			buffer.WriteString(" 배당금: ")
			buffer.WriteString(strconv.Itoa(share))
			buffer.WriteString("\n")
		}
	}
	fmt.Printf("- 수익분배조회:\n%s\n", buffer.String())

	return buffer.String(), nil
}

// 만기 지분 수거 : musicId
func (t *MusicAsset) expire(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	// musicID로 조회
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}
	err := json.Unmarshal(musicAsBytes, &music)

	list := music.Contracts
	if list == nil {
		return "the list is empty !!", nil
	}

	now := time.Now().Format("2006-01-02")

	now2int, _ := strconv.Atoi(now)

	// 만기된 지분 수거
	for _, l := range list {
		expired2int, _ := strconv.Atoi(l.Expired)
		left := expired2int - now2int
		if left < 0 {
			music.EntireStake += l.Stake
		}
	}
	// 원장에 저장
	musicAsBytes, _ = json.Marshal(music)
	err = APIstub.PutState(args[0], musicAsBytes)
	if err != nil {
		return "Can't update the whole stake !!", err
	}

	return "", nil
}

func main() {
	err := shim.Start(new(MusicAsset)) // shim에 chaincode 객체를 인자로 전달
	if err != nil {
		fmt.Printf("error creating new Smart Contract: %s", err)
	}
}
