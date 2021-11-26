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
	}else if fn == "fill" {
		result, err = t.fillInContract(APIstub, arg)
	} else if fn == "query" {
		result, err = t.queryContract(APIstub, arg)
	} else if fn == "expire" {
    result, err = t.expire(APIstub, arg)
  } else if fn == "share" {
		result, err = t.shareProfit(APIstub, arg)
	} else{
		return shim.Error("Not supporte chaincode function!!")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

// 음원 등록: 제목, 아티스트, 음악길이, 지분 생성 개수를 입력받음 
func (t *MusicAsset) registerMusic(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}
	// 원작자가 요청한 개수만큼 지분 생성
	value, _ := strconv.Atoi(args[3])
	points := t.IssueStock(APIstub, value)

	var music = Music{Name: args[0], Artist: args[1], Length: args[2], Points: points, Authorized: nil}
	musicAsBytes, _ := json.Marshal(music)
	err := APIstub.PutState(args[0], musicAsBytes)

	if err != nil {
		return "", fmt.Errorf("failed to set asset: %s", err)
	}

	return string(musicAsBytes), nil
}

// 계약서 생성 : musicId, artist, myId, stake, "date", duration
// date 추가하기
func (t *MusicAsset) setContract(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 5 {
		return "", fmt.Errorf("Incorrect number of arguments !!")
	}
	// 음원명에 해당하는 Music 갖고오기 ?
	musicAsBytes, err := APIstub.GetState(args[1])
	music := Music{}
	json.Unmarshal(musicAsBytes, &music)
  
	// 만료기한 설정
	duration, _ := strconv.Atoi(args[4])
	now := time.Now()
	convDays := duration
	expired := now.AddDate(0, 0, convDays).Format("2006-01-02")
	contract := Contract{MusicID: args[0], Artist: args[1], Buyer: args[2], Stake: args[3], Expired: expired}

	music.Contracts = Append(music.Contracts, contract)
	musicAsBytes, _ = json.Marshal(music)

	err = APIstub.PutState(args[0], musicAsBytes)
	return "the contract has successfully set", err
}

// 이용자 계약서 내용 채우기  : 음원ID -> 계약내용 추가
// 후원자 음원 후원         - stake N 
// 2차 창작자 음원 대여 구매 - stake 0 
// musicID, entireStake, Buyer, Stake, Expired, Date(자동생성)
func (t *MusicAsset) fillInContract(APIstub shim.ChaincodeStubInterface, args []string)(string, error){
  if len(args) != 5 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}

  // musicID로 조회
	musicAsBytes, _ := APIstub.GetState(args[0])
	music := Music{}

	err := json.Unmarshal(musicAsBytes, &music)

  list := music.Contracts
  preset := list[0]

  // 후원자는 entireStake 줄이기
  if args[3] != 0{
    // 남은 stake 체크
    if music.EntireStake - args[1] < 0{
      return "", fmt.Errorf("No more stake. Out of EntireStake.")
    }
    
    music.EntireStake -= args[1]
  }

  // 계약내용 채우기
  now := time.Now()
	date := now.Format("2006-01-02")

  contract := Contract{MusicID: preset[0], Artist: preset[1], Buyer: args[2], Stake: args[3], Date: date, Expired: args[4]}

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
			buffer.WriteString(d.Stake)
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
	whole := args[1] // 총수익
	err = json.Unmarshal(musicAsBytes, &music)

	list := music.Contracts
	if list == nil {
		return "the list is empty !!", nil
	}
	list_d := Contracts{} // 유효 후원자 list

	wholeStake := music.EntireStake
	time := time.Now().Format("2006-01-02")

	// 지분
	for _, l := range list {
		if l.Stake != 0 && (strconv.Atoi(l.Expired)-strconv.Atoi(time) >= 0) {
			list_d = append(list_d, l)
			wholeStake += l.Stake
		}
	}

	var buffer bytes.Buffer

	if list_d != nil {
		buffer.WriteString(music.MusicID)
		buffer.WriteString(" 수익분배표\n")
		for _, d := range list_d {
			const share = whole * d.stake / wholeStake
			buffer.WriteString("소유주: ")
			buffer.writeString(d.Owner)
			buffer.WriteString(" 배당금: ")
			buffer.WriteString(strconv.Itoa(wholeStake))
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

	time := time.Now().Format("2006-01-02")
	// 만기된 지분 수거
	for _, l := range list {
		if strconv.Atoi(l.Expired)-strconv.Atoi(time) < 0 {
			music.EntireStake += l.Stake
		}
	}
	// 원장에 저장
	musicAsBytes, _ = json.Marshal(music)

	return buffer.String(), nil
}

func main() {
	err := shim.Start(new(MusicAsset)) // shim에 chaincode 객체를 인자로 전달
	if err != nil {
		fmt.Printf("error creating new Smart Contract: %s", err)
	}
}