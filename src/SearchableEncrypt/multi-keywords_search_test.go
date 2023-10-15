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
	MatrixPrint(&sk.M1)
	MatrixPrint(&sk.M2)
	MatrixPrint(&sk.S)
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
	//keywords := ReadKeyWords()
	//sk = SetUp(len(keywords))

	sk = ReadSk()
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
	//docInx1 := BuildIndexForNode(document1, sk)
	//docInx2 := BuildIndexForNode(document2, sk)

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

func TestQueryWithSplitMat(t *testing.T) {
	var sk *SecretKey
	keywords := ReadKeyWords()
	sk = SetUp(len(keywords))

	//sk = ReadSk()
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
	//docInx1 := BuildIndex(document1, sk)
	//docInx2 := BuildIndex(document2, sk)
	docInx1 := BuildIndexForNode(document1, sk)
	docInx2 := BuildIndexForNode(document2, sk)

	mat1 := splitMat(&sk.M1)
	mat2 := splitMat(&sk.M2)

	docVec1 := RestoreDocumentVecFromDocument(docInx1, mat1, mat2)
	docVec2 := RestoreDocumentVecFromDocument(docInx2, mat1, mat2)

	documents := []DocumentVec{
		*docVec1, *docVec2,
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
		result_q1 := QueryForUserWithSplitMat(query1, documents, sk)
		println(result_q1[0] > result_q1[1])
		result_q2 := QueryForUserWithSplitMat(query2, documents, sk)
		println(result_q2[0] < result_q2[1])
	}
}

func TestWithComparing(t *testing.T) {
	var sk *SecretKey
	sk = ReadSk()
	document1 := []string{
		"肿瘤科",
		"四肢乏力",
		"记忆力衰退",
	}

	docInx1 := BuildIndexForNode(document1, sk)

	mat1 := splitMat(&sk.M1)
	mat2 := splitMat(&sk.M2)

	docVec1 := RestoreDocumentVecFromDocument(docInx1, mat1, mat2)

	documents := []DocumentVec{
		*docVec1,
	}

	query1 := []string{
		"肿瘤科",
		"四肢乏力",
		"记忆力衰退",
	}

	result_q1 := QueryForUserWithSplitMat(query1, documents, sk)
	println(result_q1[0])

	//doctest := GenerateDocumentVector(document1)
	//p1, p2 := splitVec(doctest, &sk.S, 0)
	//p1.MulVec(sk.M1.T(), p1)
	//p2.MulVec(sk.M2.T(), p2)
	//door := TrapDoorWithSplitMat(query1, sk)
	//var result1 mat.VecDense
	//var result2 mat.VecDense
	//
	//result1.MulVec(p1.T(), &door.Q1)
	//fmt.Println(result1.At(0, 0))
	//result2.MulVec(p2.T(), &door.Q2)
	//fmt.Println(result2.At(0, 0))

	docInx_1 := BuildIndex(document1, sk)
	docs := []Document{
		*docInx_1,
	}

	result_q1 = QueryForUser(query1, docs, sk)
	println(result_q1[0])

}

func TestGenerateSk(t *testing.T) {
	//GenerateSk()
	sk := ReadSk()
	MatrixPrint(&sk.M1)
}
