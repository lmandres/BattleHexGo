package helper

import (
    "math"
)

const FirstPlayer = 1
const SecondPlayer = 2
const VerticalBoard = 4
const HorizontalBoard = 8
const DiamondBoard = 16

func GetPlayerSuits(player int) string {

    var returnSuits string

    if player == FirstPlayer {
        returnSuits = "s,d"
    } else if player == SecondPlayer {
        returnSuits = "c,h"
    }

    return returnSuits
}

func GetOpponentSuits(player int) string {

    var returnSuits string

    if player == FirstPlayer {
        returnSuits = "c,h"
    } else if player == SecondPlayer {
        returnSuits = "s,d"
    }

    return returnSuits
}

func GetCellRadius(imageWidth int, imageHeight int, boardRows int, boardColumns int, boardShape int) float64 {

    var cellRadius float64 = 0
    var cellYRadius float64 = 0

    if imageWidth > 0 && imageHeight > 0 && boardRows > 0 && boardColumns > 0 {
        switch boardShape {
            case VerticalBoard+FirstPlayer, VerticalBoard+SecondPlayer:
                cellYRadius = (float64(imageHeight)/(2+(((float64(boardRows)-1)+(float64(boardColumns)-1)+3)*(1+math.Sin(math.Pi/6)))+(2/math.Sin(math.Pi/6))))
                cellRadius = (float64(imageWidth)/(((float64(boardColumns)+1)*math.Sin(math.Pi/3))+((float64(boardRows)-1)*math.Sin(2*math.Pi/3))+(2/math.Sin(math.Pi/3))+(2*math.Cos(math.Pi/6))))
        }

        if cellRadius > cellYRadius && cellYRadius > 0 {
            cellRadius = cellYRadius
        }
    }

    return cellRadius
}

func GetStartCoords(imageWidth int, imageHeight int, boardRows int, boardColumns int, boardShape int) (float64, float64) {

    var startXCoord float64 = 0
    var startYCoord float64 = 0

    var cellRadius float64 = GetCellRadius(imageWidth, imageHeight, boardRows, boardColumns, boardShape)

    if imageWidth > 0 && imageHeight > 0 && boardRows > 0 && boardColumns > 0 {
        switch boardShape {
            case VerticalBoard+FirstPlayer, VerticalBoard+SecondPlayer:
                startXCoord = float64((float64(imageWidth)/2)-((float64(boardColumns)-float64(boardRows))*GetCellHalfWidth(cellRadius)/2))
                startYCoord = float64((float64(imageHeight)/2)-((float64(boardColumns)-float64(boardRows))*GetCellExtendedLength(cellRadius)/2))-((float64(boardRows)-1)*GetCellExtendedLength(cellRadius))
        }
    }

    return startXCoord, startYCoord
}

func GetCellHalfWidth(cellRadius float64) float64 {
    return cellRadius*math.Cos(math.Pi/6)
}

func GetCellExtendedLength(cellRadius float64) float64 {
    return cellRadius*(1+math.Sin(math.Pi/6))
}

