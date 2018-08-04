package utils

import (
	"fmt"
)

var (
	CmdSizeMap = map[int]int{
		0x00: 45,
		0x03: 17,
		0x04: 17,
		0x05: 20, // OK
		0x06: 20, // OK
		0x07: 17,
		0x08: 17,
		0x09: 35,
		0x0F: 16,
		0x14: 16, // OK
		0x15: 16, // OK
		0x16: 16, // OK
		0x21: 20, // OK
		0x2C: 29, // OK
		0x32: 53,
		0x34: 45,
		0x35: 1049,
		0x36: 16, //OK
		0x5F: 11,

		// special size command
		0x01: 0,
		0x02: 0,
		0x0C: 0,
		0x10: 0,
		0x33: 0,
		0x4B: 0,

		// standard layout command, value is offset n
		0x0A: -2,
		0x0D: -2,
		0x0E: -2,
		0x12: -2,
		0x1A: -2,
		0x1B: -2,
		0x28: -2,
		0x29: -2,
		0x2A: -2,
		0x2E: -2,
		0x2F: -2,
		0x37: -2,
		0x47: -2,
		0x48: -2,
		0x4C: -2,
		0x4E: -2,
		0x52: -2,
		0xF5: -5,
		0xF6: -5,
		0xF8: -4,
		0xF9: -2,
		0xFA: -7,
		0xFB: -2,
		0xFC: -2,
		0xFD: -7,
		0xFE: -15,
		0xFF: -34,
	}
	CmdNameMap = map[int]string{
		0x02: "set rally point",
		0x03: "start/resume research upgrade",
		0x04: "pause/cancel research upgrade",
		0x05: "start/resume unit construction",
		0x06: "pause/cancel unit construction",
		0x07: "start/resume structure construction",
		0x08: "pause/cancel structure construction",
		0x09: "place structure",
		0x0A: "sell structure",

		0x0C: "ungarrison structure (?)",
		0x0D: "attack",
		0x0E: "force-fire",

		0x10: "garrison structure",

		0x14: "move unit",
		0x15: "attack-move unit",
		0x16: "force-move unit",
		0x1A: "stop unit",

		0x21: "3s heartbeat",

		0x28: "start repair structure",
		0x29: "stop repair structure",
		0x2A: "'Q' select",
		0x2C: "formation-move preview",
		0x2E: "stance change",
		0x2F: "waypoint/planning mode (?)",
		0x37: "'scroll'",
		0x4E: "player power",
		0xF5: "drag selection box and/or select units/structures",
		0xF8: "left click",
		0xF9: "unit ungarrisons structure (automatic event) (?)",
		0xFA: "create group",
		0xFB: "select group",
	}
)

type BodyChunk struct {
	TimeCode  uint32
	ChunkType byte
	ChunkSize uint32
	Data      []byte
	Zero      uint32
}

type ReplayBody struct {
	Size   int         // Chunk count
	Chunks []BodyChunk // Chunk list
}

// 计算玩家手速
func (rb *ReplayBody) CalcAPM(players []PlayerDetail) (result []PlayerDetail) {
	playermap := make(map[int]int)
	playertime := make(map[int]int)
	// 过滤掉视角切换
	for i := 0; i < rb.Size; i++ {
		timeCode := int(rb.Chunks[i].TimeCode)
		switch rb.Chunks[i].ChunkType {
		case 1:
			chunkData := rb.Chunks[i].Data
			number_of_commands := int(ToUInt32LE(chunkData[1:5]))
			payload := chunkData[5:]
			payloadArr := make([][]byte, 0)
			if number_of_commands > 1 {
				// 多条命令，需要拆分命令
				for j, ptr, pLen := 0, 0, number_of_commands; j < pLen; j++ {
					commandId := payload[ptr]
					cmdsize := CmdSizeMap[int(commandId)]
					var arr []byte
					if cmdsize < 0 { // Variable-length
						arr = ParseVariableLen(payload[ptr:], cmdsize)
					} else if cmdsize == 0 { // Special-length
						arr = ParseSpecialLen(payload[ptr:])
					} else { // Fixed-length
						arr = ParseFixedLen(payload[ptr:], cmdsize)
					}
					payloadArr = append(payloadArr, arr)
					ptr += len(arr)
				}
			} else {
				// 一条命令
				payloadArr = append(payloadArr, payload)
			}
			// fmt.Println("parsed payload arr:", payloadArr)
			if len(payloadArr) != number_of_commands {
				fmt.Println("parsed error:", number_of_commands, payloadArr)
			}
			// 遍历 payloadArr 获取玩家操作数
			// 获取每位玩家的 finaltimecode
			for b, pLen := 0, len(payloadArr); b < pLen; b++ {
				item := payloadArr[b]
				// 过滤掉心跳 heartbeats
				if int(item[0]) == 0x21 {
					continue
				}
				playerIndex := int(item[1])/8 - 2 // playerId/8 - k, k is 2 for RA3
				playermap[playerIndex]++
				playertime[playerIndex] = timeCode
			}
		case 2: // 记录了摄像头的位置和角度，每 1/15 秒一帧
			// chunkData := rb.Chunks[i].Data
			// playerIndex := int(ToUInt32LE(chunkData[2:6]))
			// timeCode := int(ToUInt32LE(chunkData[7:11]))
			// fmt.Println("matched type 2", playerIndex, timeCode)
			fallthrough
		default:
			continue
		}
	}
	// fmt.Println("player map:", playermap)
	// fmt.Println("player time:", playertime)
	for k, v := range playermap {
		minutes := float64(playertime[k] / 15 / 60)
		players[k].Apm = int(float64(v) / minutes)
	}
	return players
}

func ParseVariableLen(payload []byte, offset int) (arr []byte) {
	offset = -offset
	// fmt.Println("parsed variable len", payload, len(payload), offset)
	end, x := offset, payload[offset]
	for {
		if x == byte(0XFF) {
			break
		}
		c := int((x >> 4) + 1)
		end += 4*c + 1
		x = payload[end]
	}
	return payload[:end+1]
}

func ParseFixedLen(payload []byte, cmdsize int) (arr []byte) {
	return payload[:cmdsize]
}

func ParseSpecialLen(payload []byte) (arr []byte) {
	end := 0
	switch int(payload[0]) {
	case 0x01:
		if int(payload[2]) == 0xFF {
			end += 3
		} else if int(payload[7]) == 0xFF {
			end += 8
		} else {
			l := int(payload[17]) + 1
			end += 4*l + 32
		}
	case 0x02:
		l := (int(payload[24])+1)*2 + 26
		end += l
	case 0x0C:
		l := int(payload[3]) + 1
		end += 4*l + 5
	case 0x10:
		l := 0
		if int(payload[2]) == 0x14 {
			l = 12
		} else if int(payload[2]) == 0x04 {
			l = 13
		} else {
			l = 99999
		}
		end += l
	case 0x33:
		l := ParseUuid(payload)
		end += l
	case 0x4B:
		l := 0
		if int(payload[2]) == 0x04 {
			l = 8
		} else if int(payload[2]) == 0x07 {
			l = 16
		} else {
			l = 99999
		}
		end += l
	default: // 未知code
		for end, pLen := 0, len(payload); end < pLen; end++ {
			if int(payload[end]) != 0xFF {
				break
			}
		}
		if int(payload[end]) != 0xFF {
			end = 99999
		} else {
			end++
		}
	}
	if end >= 99999 {
		return payload
	}
	return payload[:end]
}

func ParseUuid(payload []byte) (end int) {
	l := int(payload[3])
	if l > len(payload) {
		return 0
	}
	end = l + 5
	l = int(payload[end])
	if l > len(payload) {
		return 0
	}
	end += 2*l + 2
	l = int(payload[end])
	if l > len(payload) {
		return 0
	}
	end += 5
	return end
}
