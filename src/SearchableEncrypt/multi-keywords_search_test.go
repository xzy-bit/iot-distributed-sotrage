package SearchableEncrypt

import (
	"gonum.org/v1/gonum/mat"
	"log"
	"testing"
)

func TestMatrixPrint(t *testing.T) {
	v := mat.NewVecDense(4, []float64{1, 2, 3, 4})
	MatrixPrint(v)
}

func TestGenerateInvertibleMatrix(t *testing.T) {
	var Inv mat.Dense
	matrix := GenerateRandomMatrix(20)
	err := Inv.Inverse(matrix)
	if err != nil {
		log.Fatal(err)
	}
	MatrixPrint(matrix)
	println()
	MatrixPrint(&Inv)

	var result mat.Dense
	result.Mul(matrix, &Inv)
	MatrixPrint(&result)
}

func TestSetUp(t *testing.T) {
	sk := SetUp(10)
	MatrixPrint(sk.M1)
	MatrixPrint(sk.M2)
	MatrixPrint(sk.S)
}

//func TestAddKeyWords(t *testing.T) {
//	strs := []string{
//		"精神科",
//		"五官科",
//		"口腔科",
//		"骨科",
//		"咽喉科",
//		"胃科",
//		"肿瘤科",
//		"妇科",
//		"心率不齐",
//		"食欲不振",
//		"四肢乏力",
//		"记忆力衰退",
//	}
//	AddKeyWords(strs)
//}

func TestQuery(t *testing.T) {
	var sk *SecretKey
	keywords := ReadKeyWords()
	sk = SetUp(len(keywords))

	document1 := []string{
		"肿瘤科",
		"四肢乏力",
		"记忆力衰退",
	}
	document2 := []string{
		"肿瘤科",
		"四肢乏力",
		"心率不齐",
	}
	docInx1 := BuildIndex(document1, sk)
	docInx2 := BuildIndex(document2, sk)
	documents := []Document{
		*docInx1, *docInx2,
	}

	query1 := []string{
		"肿瘤科",
		"四肢乏力",
		"记忆力衰退",
	}

	query2 := []string{
		"肿瘤科",
		"四肢乏力",
		"心率不齐",
	}

	for i := 0; i < 10; i++ {
		result_q1 := QueryForUser(query1, documents, sk)
		println(result_q1[0] > result_q1[1])
		result_q2 := QueryForUser(query2, documents, sk)
		println(result_q2[0] < result_q2[1])
	}
}
