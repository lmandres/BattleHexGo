package frontend

import (
    "math"
)

const (
    FirstPlayer = 1
    SecondPlayer = 2
    VerticalBoard = 4
    HorizontalBoard = 8
    DiamondBoard = 16
)

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

func GetBoardLabel(rowColumnIn int) string {
    var indexValues [13]string = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
    return indexValues[rowColumnIn-1]
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

func GetObtuseBorderXCoord(cellRadius float64, boardColumns int) float64 {
    return ((float64(boardColumns)-1)*GetCellHalfWidth(cellRadius))
}

func GetObtuseBorderYCoord(cellRadius float64, boardColumns int) float64 {
    return ((float64(boardColumns)-1)*GetCellExtendedLength(cellRadius))
}

func GetAcuteBorderXCoord(cellRadius float64, boardRows int) float64 {
    return 0
}

func GetAcuteBorderYCoord(cellRadius float64, boardRows int) float64 {
    return (float64(boardRows)-1)*(2*GetCellExtendedLength(cellRadius))
}

func GetAcuteBorderLength(cellRadius float64) float64 {
    return (cellRadius/math.Sin(math.Pi/6))+(2*(cellRadius*(1+math.Sin(math.Pi/6))))
}

func GetObtuseBorderLength(cellRadius float64) float64 {
    return (cellRadius/math.Sin(math.Pi/3))+(2*(cellRadius*(math.Cos(math.Pi/6))))
}

func GetCellXCoord(cellRadius float64, startXCoord float64, boardRows int, boardColumns int, rowIn int, columnIn int, boardShape int) float64 {

    var returnXCoord float64 = 0

    switch boardShape {
        case VerticalBoard+FirstPlayer:
            returnXCoord = startXCoord+(float64(columnIn-rowIn)*GetCellHalfWidth(cellRadius))
        case VerticalBoard+SecondPlayer:
            returnXCoord = startXCoord+((float64(boardColumns-columnIn+1)-float64(boardRows-rowIn+1))*GetCellHalfWidth(cellRadius))
    }

    return returnXCoord
}

func GetCellYCoord(cellRadius float64, startYCoord float64, boardRows int, boardColumns int, rowIn int, columnIn int, boardShape int) float64 {

    var returnYCoord float64 = 0

    switch boardShape {
        case VerticalBoard+FirstPlayer:
            returnYCoord = startYCoord+((float64(columnIn-rowIn)*GetCellExtendedLength(cellRadius))+(float64(rowIn-1)*(2*GetCellExtendedLength(cellRadius))))
        case VerticalBoard+SecondPlayer:
            returnYCoord = startYCoord+((float64(boardColumns-columnIn+1)-float64(boardRows-rowIn+1))*GetCellExtendedLength(cellRadius))+((float64(boardRows-rowIn))*(2*GetCellExtendedLength(cellRadius)))
    }

    return returnYCoord
}

func GetPointXCoord(cellRadius float64, cellXCoordIn float64, pointIndexIn int, boardShape int) float64 {

    var returnXCoord float64 = 0

    switch boardShape {
        case VerticalBoard+FirstPlayer, VerticalBoard+SecondPlayer:
            returnXCoord = cellXCoordIn+(cellRadius*math.Cos((float64(pointIndexIn)*math.Pi/3)+math.Pi/6))
    }

    return returnXCoord
}

func GetPointYCoord(cellRadius float64, cellYCoordIn float64, pointIndexIn int, boardShape int) float64 {

    var returnYCoord float64 = 0

    switch boardShape {
        case VerticalBoard+FirstPlayer, VerticalBoard+SecondPlayer:
            returnYCoord = cellYCoordIn+(cellRadius*math.Sin((float64(pointIndexIn)*math.Pi/3)+math.Pi/6))
    }

    return returnYCoord
}
