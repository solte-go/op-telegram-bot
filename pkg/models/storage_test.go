package models

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type testCase struct {
	name           string
	data           string
	expectedResult error
}

type storageModelTests struct {
	suite.Suite
	testCases []testCase
	parser    Parser
}

func (suite *storageModelTests) SetupSuite() {
	suite.testCases = []testCase{
		{
			name:           "Success",
			data:           "Ohjelmoija, Programmer,Программист, Profession;",
			expectedResult: nil,
		}, {
			name:           "Bad line syntax",
			data:           "Ohjelmoija, Programmer,Программист, Profession",
			expectedResult: ErrBadDataInLine,
		},
		{
			name:           "Missig args",
			data:           "Ohjelmoija, ,Программист, Profession",
			expectedResult: ErrBadDataInLine,
		},
	}

	suite.parser = NewParser()
}

//func (suite *storageModelTests) Test_Parse() {
//	for _, ts := range suite.testCases {
//		suite.Run(ts.name, func() {
//			//w := Words{}
//			_, err := suite.parser.Parse(ts.data)
//			suite.NoError(err)
//			suite.Equal(ts.expectedResult, err)
//		})
//	}
//}

func (suite *storageModelTests) Test_CheckData() {
	for _, ts := range suite.testCases {
		suite.Run(ts.name, func() {
			_, err := suite.parser.Parse(ts.data)
			if err != nil {
				if ts.name == "Success" {
					suite.NoError(err)
				} else {
					suite.Error(ts.expectedResult, "Error: ", err)
				}
			}
		})
	}
}

func (suite *storageModelTests) Test_ParseData() {
	data := "Ohjelmoija, Programmer,Программист, Profession;"

	w, err := suite.parser.Parse(data)
	suite.NoError(err)
	suite.Equal("o", w[0].Letter)
	suite.Equal("ohjelmoija", w[0].Suomi)
	suite.Equal("programmer", w[0].English)
	suite.Equal("программист", w[0].Russian)
	suite.Equal("profession", w[0].Topic)
}

func TestStorageModelSuite(t *testing.T) {
	suite.Run(t, new(storageModelTests))
}
