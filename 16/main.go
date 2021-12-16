package main

import (
	"fmt"
	"log"
	"strconv"
)

const PUZZLE_INPUT = "805311100469800804A3E488ACC0B10055D8009548874F65665AD42F60073E7338E7E5C538D820114AEA1A19927797976F8F43CD7354D66747B3005B401397C6CBA2FCEEE7AACDECC017938B3F802E000854488F70FC401F8BD09E199005B3600BCBFEEE12FFBB84FC8466B515E92B79B1003C797AEBAF53917E99FF2E953D0D284359CA0CB80193D12B3005B4017968D77EB224B46BBF591E7BEBD2FA00100622B4ED64773D0CF7816600B68020000874718E715C0010D8AF1E61CC946FB99FC2C20098275EBC0109FA14CAEDC20EB8033389531AAB14C72162492DE33AE0118012C05EEB801C0054F880102007A01192C040E100ED20035DA8018402BE20099A0020CB801AE0049801E800DD10021E4002DC7D30046C0160004323E42C8EA200DC5A87D06250C50015097FB2CFC93A101006F532EB600849634912799EF7BF609270D0802B59876F004246941091A5040402C9BD4DF654967BFDE4A6432769CED4EC3C4F04C000A895B8E98013246A6016CB3CCC94C9144A03CFAB9002033E7B24A24016DD802933AFAE48EAA3335A632013BC401D8850863A8803D1C61447A00042E3647B83F313674009E6533E158C3351F94C9902803D35C869865D564690103004E74CB001F39BEFFAAD37DFF558C012D005A5A9E851D25F76DD88A5F4BC600ACB6E1322B004E5FE1F2FF0E3005EC017969EB7AE4D1A53D07B918C0B1802F088B2C810326215CCBB6BC140C0149EE87780233E0D298C33B008C52763C9C94BF8DC886504E1ECD4E75C7E4EA00284180371362C44320043E2EC258F24008747785D10C001039F80644F201217401500043A2244B8D200085C3F8690BA78F08018394079A7A996D200806647A49E249C675C0802609D66B004658BA7F1562500366279CCBEB2600ACCA6D802C00085C658BD1DC401A8EB136100"

const TYPE_ID_LITERAL_VALUE = 4
const TYPE_ID_SUM = 0
const TYPE_ID_PRODUCT = 1
const TYPE_ID_MIN = 2
const TYPE_ID_MAX = 3
const TYPE_ID_GT = 5
const TYPE_ID_LT = 6
const TYPE_ID_EQ = 7

func hexCharToBits(c string) []bool {
	val, err := strconv.ParseUint(string(c), 16, 4)
	if err != nil {
		log.Fatalf("Failed parsing the hex characters")
	}
	bits := make([]bool, 4)
	for i := 0; i < 4; i++ {
		bits[3-i] = val&0x1 == 0x1
		val = val >> 1
	}
	return bits
}

func printBits(bits []bool) {
	for _, b := range bits {
		if b {
			print("1")
		} else {
			print("0")
		}
	}
}

func bitsToInt(bits []bool) int {
	output := 0
	for i := 0; i < len(bits); i++ {
		output <<= 1
		if bits[i] {
			output += 1
		}
	}
	return output
}

type packet struct {
	version int
	typeId  int
	payload []interface{}
}

func (p *packet) getPayloadPackets() []packet {
	output := make([]packet, len(p.payload))
	for i, val := range p.payload {
		output[i] = val.(packet)
	}
	return output
}

func (p *packet) getPayloadValues() []int {
	output := make([]int, len(p.payload))
	for i, pl := range p.getPayloadPackets() {
		output[i] = pl.evaluate()
	}
	return output
}

func (p *packet) evaluate() int {
	switch p.typeId {
	case TYPE_ID_LITERAL_VALUE:
		return p.payload[0].(int)
	case TYPE_ID_SUM:
		output := 0
		for _, val := range p.getPayloadValues() {
			output += val
		}
		return output
	case TYPE_ID_PRODUCT:
		output := 1
		for _, val := range p.getPayloadValues() {
			output *= val
		}
		return output
	case TYPE_ID_MIN:
		output := 2147483647
		for _, val := range p.getPayloadValues() {
			if val < output {
				output = val
			}
		}
		return output
	case TYPE_ID_MAX:
		output := 0
		for _, val := range p.getPayloadValues() {
			if val > output {
				output = val
			}
		}
		return output
	case TYPE_ID_GT:
		values := p.getPayloadValues()
		if values[0] > values[1] {
			return 1
		} else {
			return 0
		}
	case TYPE_ID_LT:
		values := p.getPayloadValues()
		if values[0] < values[1] {
			return 1
		} else {
			return 0
		}
	case TYPE_ID_EQ:
		values := p.getPayloadValues()
		if values[0] == values[1] {
			return 1
		} else {
			return 0
		}
	default:
		log.Fatalf("Unknown operation")
		return -1 // just to make compiler happy
	}
}

type bitstream struct {
	bits []bool
}

func (b *bitstream) consumeBit() bool {
	val := b.bits[0]
	b.bits = b.bits[1:]
	return val
}

func (b *bitstream) consumeBits(len int) []bool {
	consumed, rest := b.bits[:len], b.bits[len:]
	b.bits = rest
	return consumed
}

func (b *bitstream) consumeInt(len int) int {
	return bitsToInt(b.consumeBits(len))
}

func getPacketTree(b *bitstream) (packet, int, int) {

	payload := make([]interface{}, 0)
	pkt := packet{
		version: b.consumeInt(3),
		typeId:  b.consumeInt(3),
		payload: payload,
	}

	versionSum := pkt.version
	consumedBits := 6

	if pkt.typeId == TYPE_ID_LITERAL_VALUE {
		consumedSubBits := 0
		valueBits := []bool{}
		for {
			continuationBit := b.consumeBit()
			chunk := b.consumeBits(4)
			consumedSubBits += 5
			valueBits = append(valueBits, chunk...)
			if !continuationBit {
				break
			}
		}
		literalValue := bitsToInt(valueBits)
		pkt.payload = append(pkt.payload, literalValue)
		consumedBits += consumedSubBits
	} else {
		lengthBit := b.consumeBit()
		consumedSubBits := 1
		if lengthBit {
			nrSubPackets := b.consumeInt(11)
			consumedSubBits += 11
			for i := 0; i < nrSubPackets; i++ {
				subPkt, subSum, subCnt := getPacketTree(b)
				consumedSubBits += subCnt
				versionSum += subSum
				pkt.payload = append(pkt.payload, subPkt)
			}
		} else {
			subPacketLen := b.consumeInt(15)
			consumedSubBits += 15
			for {
				subPkt, subSum, subCnt := getPacketTree(b)
				consumedSubBits += subCnt
				versionSum += subSum
				pkt.payload = append(pkt.payload, subPkt)

				subPacketLen -= subCnt
				if subPacketLen < 0 {
					log.Fatalf("???")
				}
				if subPacketLen == 0 {
					break
				}
			}

		}
		consumedBits += consumedSubBits
	}
	return pkt, versionSum, consumedBits
}

func parseInput(input string) {
	bits := []bool{}
	for _, c := range input {
		curBits := hexCharToBits(string(c))
		bits = append(bits, curBits...)
	}

	b := bitstream{bits: bits}

	pkt, versionSum, _ := getPacketTree(&b)
	print(fmt.Sprintf("Total version sum: %d\n", versionSum))
	print(fmt.Sprintf("Result: %d\n", pkt.evaluate()))
}

func main() {
	parseInput(PUZZLE_INPUT)
}
