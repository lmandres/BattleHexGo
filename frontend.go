package main

import (
    "fmt"
    "net/http"

    "battle-hex-go/helper"
)

func BattleHexJSHandler(w http.ResponseWriter, r *http.Request) {

    var cardset string = "/static/cardstux"

    var player int = helper.FirstPlayer
    var imageWidth int = 320
    var imageHeight int = 550
    var boardRows int = 13
    var boardColumns int = 13
    var boardShape int = helper.VerticalBoard

    var xCoordStr string = ""
    var yCoordStr string = ""

    var opponentMatch string
    var playerMatch string

    if player == helper.FirstPlayer {
        opponentMatch = "Opposite"
        playerMatch = "Same"
    } else if player == helper.SecondPlayer {
        opponentMatch = "Same"
        playerMatch = "Opposite"
    }

    startXCoord, startYCoord := helper.GetStartCoords(imageWidth, imageHeight, boardRows, boardColumns, boardShape+player)
    cellRadius := helper.GetCellRadius(imageWidth, imageHeight, boardRows, boardColumns, boardShape+player)

    pagePartStart := `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<!-- The HTML 4.01 Transitional DOCTYPE declaration-->
<!-- above set at the top of the file will set     -->
<!-- the browser's rendering engine into           -->
<!-- "Quirks Mode". Replacing this declaration     -->
<!-- with a "Standards Mode" doctype is supported, -->
<!-- but may lead to some differences in layout.   -->

<html>
`

    pagePartHeadOpen := `
  <head>
    <meta http-equiv="content-type" content="text/html; charset=ISO-8859-1">
    <title>Battle Hex Go</title>
    <link rel="shortcut icon" href="/static/favicon.ico" />
    <style>
        .loaderanimation {
            position: absolute;
            left: 50%;
            top: 50%;
            z-index: 10;
            margin: -35px 0 0 -35px;
            border: 8px solid #f3f3f3;
            border-radius: 50%;
            border-top: 8px solid #3498db;
            width: 50px;
            height: 50px;
            animation: spin 2s linear infinite;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
    <script type="text/javascript" src="/static/numeric-1.2.6.min.js"></script>
    <script type="text/javascript" src="/static/math.min.js"></script>
`

    pagePartJSGameCode := `
    <script type="text/javascript">

        let cardset = "%s";
    
        let player = %d;
        let computer = (player %% 2) + 1;

        let svgNS = "http://www.w3.org/2000/svg";
        let svgID = "svgBoard";
        
        let indexValues = ["a", "2", "3", "4", "5", "6", "7", "8", "9", "10", "j", "q", "k"];

        let playerFlipCard = null;
        let playerPlayCard = null;

        let computerFlipCard = null;
        let computerPlayCard = null;

        let playerSuits = "%s".split(",");
        let computerSuits = "%s".split(",");
        let deleteLine = null;

        let gameBoard = [];
        let voltageBoardIndexes = [];

        let testVoltage = 100;
        let randomMoveLimit = 100;

        let loaderTimeout = 50;
        
        function supportsSVG() {
            return !!document.createElementNS && !!document.createElementNS(svgNS, "svg").createSVGRect;
        }
        
        function getXCoord(row, column) {
            return %s;
        }
        
        function getYCoord(row, column) {
            return %s;
        }
        
        function makeCardMove(card, playerIn) {

            let gameBoardSet = true;

            if (gameBoard.length != indexValues.length+2) {
                gameBoardSet = false;
            } else {
                for (let rowIndex = 0; (rowIndex < indexValues.length+2) || !gameBoardSet; rowIndex++) {
                    if (gameBoard[rowIndex].length != indexValues.length+2) {
                        gameBoardSet = false;
                    }
                }
            }

            if (gameBoardSet) {

                if (playerIn == player) {
            
                    if (
                        playerFlipCard == null ||
                        playerPlayCard != null ||
                        (
                            playerPlayCard == null &&
                            card.substr(0,1).toLowerCase() != playerFlipCard.substr(0,1).toLowerCase()  
                        )
                    ) {

                        if (playerFlipCard == null) {
                            playerFlipCard = card;
                        } else if (playerPlayCard == null) {

                            let playerMove = getPlayerRowColumn(playerFlipCard, card);

                            if (
                                playerMove["row"] &&
                                playerMove["column"] &&
                                gameBoard[playerMove["row"]][playerMove["column"]] == null
                            ) {
                                playerPlayCard = card;
                            }
            
                        } else {

                            let playerMove = getPlayerRowColumn(playerFlipCard, playerPlayCard);

                            drawHexCell(playerMove["row"], playerMove["column"], "makeCellMove("+playerMove["row"]+","+playerMove["column"]+",player);");

                            playerFlipCard = card;
                            playerPlayCard = null;
                        }

                        displayMove();
                    }

                } else {

                    if (computerFlipCard == null) {
                        computerFlipCard = card;
                    } else if (computerPlayCard == null) {
                        computerPlayCard = card;                    
                    } else {
                        computerFlipCard = card;
                        computerPlayCard = null;
                    }
                }
                    

                if ((playerPlayCard != null) || (playerFlipCard != null)) {
                    document.getElementById("opponentPlayCardImg").src = cardset + "/back.png";
                    document.getElementById("opponentFlipCardImg").src = cardset + "/back.png";
                }
            }
        }
        
        function makeCellMove(row, column, playerIn) {

            if (!evaluateWin(player, gameBoard) && !evaluateWin(computer, gameBoard)) {

                if (playerIn == player) {

                    if ((playerPlayCard == null) && (playerFlipCard != null)) {
                        if (deleteLine.substr(0,1) == "r") {
                            makeCardMove("b"+indexValues[column-1], playerIn);
                        } else if (deleteLine.substr(0,1) == "c") {
                            makeCardMove("r"+indexValues[row-1], playerIn);
                        }
                    }
                }

                makeCardMove("r"+indexValues[row-1], playerIn);
                makeCardMove("b"+indexValues[column-1], playerIn);

                randomizeMove(playerIn);
            }
        }
        
        function randomizeMove(playerIn) {

            if (playerIn == player) {

                if (playerFlipCard != null && playerPlayCard != null) {
                    if (Math.random() < 0.5) {
                        var tempCard = playerFlipCard;
                        playerFlipCard = playerPlayCard;
                        playerPlayCard = tempCard;
                    }

                    document.getElementById("playerFlipCardImg").src = cardset + "/" + playerSuits["b,r".split(",").indexOf(playerFlipCard.substr(0, 1))] + playerFlipCard.substr(1) + ".png";
                    document.getElementById("playerPlayCardImg").src = cardset + "/" + playerSuits["b,r".split(",").indexOf(playerPlayCard.substr(0, 1))] + playerPlayCard.substr(1) + ".png";
                }

            } else {
                if (Math.random() < 0.5) {
                    var tempCard = computerFlipCard;
                    computerFlipCard = computerPlayCard;
                    computerPlayCard = tempCard;
                }
            }
        }

        function showMoves() {

            if (
                playerFlipCard != null &&
                playerPlayCard != null &&
                computerFlipCard != null &&
                computerPlayCard != null
            ) {

                let playerMove = getPlayerRowColumn(playerFlipCard, playerPlayCard);
                let computerMove = getPlayerRowColumn(computerFlipCard, computerPlayCard);

                document.getElementById("opponentFlipCardImg").src = cardset + "/" + computerSuits["b,r".split(",").indexOf(computerFlipCard.substr(0, 1))] + computerFlipCard.substr(1) + ".png";
                document.getElementById("opponentPlayCardImg").src = cardset + "/" + computerSuits["b,r".split(",").indexOf(computerPlayCard.substr(0, 1))] + computerPlayCard.substr(1) + ".png";

                if (playerMove["row"] == computerMove["row"] && playerMove["column"] == computerMove["column"]) {
                    tieBreaker();
                } else {

                    if (player == 1) {
                        gameBoard[playerMove["row"]][playerMove["column"]] = player;
                        drawHexPiece(playerMove["row"], playerMove["column"], "red");
                        gameBoard[computerMove["row"]][computerMove["column"]] = computer;
                        drawHexPiece(computerMove["row"], computerMove["column"], "black");
                    } else {
                        gameBoard[playerMove["row"]][playerMove["column"]] = player;
                        drawHexPiece(playerMove["row"], playerMove["column"], "black");
                        gameBoard[computerMove["row"]][computerMove["column"]] = computer;
                        drawHexPiece(computerMove["row"], computerMove["column"], "red");
                    }

                    playerFlipCard = null;
                    playerPlayCard = null;
                    computerFlipCard = null;
                    computerPlayCard = null;

                    if (!evaluateWin(player, gameBoard) && !evaluateWin(computer, gameBoard)) {

                        document.getElementById("loader").classList.add("loaderanimation");
                        setTimeout(
                            function() {
                                computerMove = calculateComputerMove(computer, gameBoard);
                                makeCellMove(computerMove["row"], computerMove["column"], computer);
                                showMoves();

                                document.getElementById("loader").classList.remove("loaderanimation");
                            },
                            loaderTimeout
                        );

                    }
                }

                if (evaluateWin(player, gameBoard)) {
                    alert("Player wins!");
                } else if (evaluateWin(computer, gameBoard)) {
                    alert("Computer wins!");
                }

            } else {

                let playerCount = 0;
                let computerCount = 0;

                for (let rowIndex = 1; rowIndex <= indexValues.length; rowIndex++) {
                    for (let columnIndex = 1; columnIndex <= indexValues.length; columnIndex++) {
                        if (gameBoard[rowIndex][columnIndex] == player) {
                            playerCount++;
                        } else if (gameBoard[rowIndex][columnIndex] == computer) {
                            computerCount++;
                        }
                    }
                }

                if (
                    (playerCount < computerCount) &&
                    playerFlipCard != null &&
                    playerPlayCard != null
                ) {

                    let playerMove = getPlayerRowColumn(playerFlipCard, playerPlayCard);

                    if (computerFlipCard.substr(0, 1) == "r") {
                        makeCardMove("r" + indexValues[playerMove["row"]-1], player);
                        makeCardMove("b" + indexValues[playerMove["column"]-1], player);
                    } else {
                        makeCardMove("b" + indexValues[playerMove["column"]-1], player);
                        makeCardMove("r" + indexValues[playerMove["row"]-1], player);
                    }

                    if (player == 1) {
                        gameBoard[playerMove["row"]][playerMove["column"]] = player;
                        drawHexPiece(playerMove["row"], playerMove["column"], "red");
                    } else {
                        gameBoard[playerMove["row"]][playerMove["column"]] = player;
                        drawHexPiece(playerMove["row"], playerMove["column"], "black");
                    }

                    playerFlipCard = null;
                    playerPlayCard = null;
                    computerFlipCard = null;
                    computerPlayCard = null;

                    document.getElementById("loader").classList.add("loaderanimation");
                    setTimeout(
                        function() {
                            computerMove = calculateComputerMove(computer, gameBoard);
                            makeCellMove(computerMove["row"], computerMove["column"], computer);
                            showMoves();

                            document.getElementById("loader").classList.remove("loaderanimation");
                        },
                        loaderTimeout
                    );
                }
            }
        }

        function tieBreaker() {

            if (
                playerFlipCard != null &&
                playerPlayCard != null &&
                computerFlipCard != null &&
                computerPlayCard != null
            ) {

                let playerMove = getPlayerRowColumn(playerFlipCard, playerPlayCard);

                if (playerFlipCard == computerFlipCard) {

                    if (player == 1) {

                        let computerMove = {};

                        alert("Both move to \"" + playerFlipCard + "\", \"" + playerPlayCard + ".\"  Player wins move.");

                        gameBoard[playerMove["row"]][playerMove["column"]] = player;
                        drawHexPiece(playerMove["row"], playerMove["column"], "red");

                        if (!evaluateWin(player, gameBoard)) {

                            document.getElementById("loader").classList.add("loaderanimation");
                            setTimeout(
                                function() {
                                    computerMove = calculateComputerMove(computer, gameBoard);
                                    gameBoard[computerMove["row"]][computerMove["column"]] = computer;
                                    drawHexPiece(computerMove["row"], computerMove["column"], "black");

                                    if (playerFlipCard.substr(0, 1) == "r") {
                                        makeCardMove("r" + indexValues[computerMove["row"]-1], computer);
                                        makeCardMove("b" + indexValues[computerMove["column"]-1], computer);
                                    } else {
                                        makeCardMove("b" + indexValues[computerMove["column"]-1], computer);
                                        makeCardMove("r" + indexValues[computerMove["row"]-1], computer);
                                    }

                                    playerFlipCard = null;
                                    playerPlayCard = null;

                                    document.getElementById("opponentFlipCardImg").src = cardset + "/" + computerSuits["b,r".split(",").indexOf(computerFlipCard.substr(0, 1))] + computerFlipCard.substr(1) + ".png";
                                    document.getElementById("opponentPlayCardImg").src = cardset + "/" + computerSuits["b,r".split(",").indexOf(computerPlayCard.substr(0, 1))] + computerPlayCard.substr(1) + ".png";

                                    computerMove = calculateComputerMove(computer, gameBoard);
                                    makeCellMove(computerMove["row"], computerMove["column"], computer);

                                    document.getElementById("loader").classList.remove("loaderanimation");
                                },
                                loaderTimeout
                            );
                        }

                    } else {

                        alert("Both move to \"" + computerFlipCard + "\", \"" + computerPlayCard + ".\"  Computer wins move.");

                        let computerMove = getPlayerRowColumn(computerFlipCard, computerPlayCard);

                        playerFlipCard = null;
                        playerPlayCard = null;

                        gameBoard[computerMove["row"]][computerMove["column"]] = computer;
                        drawHexPiece(computerMove["row"], computerMove["column"], "red");
                    }
                    
                } else {

                    if (player == 2) {

                        alert("Both move to \"" + playerFlipCard + "\", \"" + playerPlayCard + ".\"  Player wins move.");

                        let computerMove = {};

                        gameBoard[playerMove["row"]][playerMove["column"]] = player;
                        drawHexPiece(playerMove["row"], playerMove["column"], "black");

                        if (!evaluateWin(player, gameBoard)) {

                            document.getElementById("loader").classList.add("loaderanimation");
                            setTimeout(
                                function() {
                                    computerMove = calculateComputerMove(computer, gameBoard);
                                    gameBoard[computerMove["row"]][computerMove["column"]] = computer;
                                    drawHexPiece(computerMove["row"], computerMove["column"], "red");

                                    if (playerFlipCard.substr(0, 1) == "r") {
                                        makeCardMove("r" + indexValues[computerMove["row"]-1], computer);
                                        makeCardMove("b" + indexValues[computerMove["column"]-1], computer);
                                    } else {
                                        makeCardMove("b" + indexValues[computerMove["column"]-1], computer);
                                        makeCardMove("r" + indexValues[computerMove["row"]-1], computer);
                                    }

                                    playerFlipCard = null;
                                    playerPlayCard = null;

                                    document.getElementById("opponentFlipCardImg").src = cardset + "/" + computerSuits["b,r".split(",").indexOf(computerFlipCard.substr(0, 1))] + computerFlipCard.substr(1) + ".png";
                                    document.getElementById("opponentPlayCardImg").src = cardset + "/" + computerSuits["b,r".split(",").indexOf(computerPlayCard.substr(0, 1))] + computerPlayCard.substr(1) + ".png";

                                    computerMove = calculateComputerMove(computer, gameBoard);
                                    makeCellMove(computerMove["row"], computerMove["column"], computer);

                                    document.getElementById("loader").classList.remove("loaderanimation");
                                },
                                loaderTimeout
                            );
                        }

                    } else {

                        alert("Both move to \"" + computerFlipCard + "\", \"" + computerPlayCard + ".\"  Computer wins move.");

                        let computerMove = getPlayerRowColumn(computerFlipCard, computerPlayCard);

                        playerFlipCard = null;
                        playerPlayCard = null;

                        gameBoard[computerMove["row"]][computerMove["column"]] = computer;
                        drawHexPiece(computerMove["row"], computerMove["column"], "black");
                    }
                }
            }
        }

        function evaluateWin(playerIn, gameBoardIn) {

            let traverseFunctions = [
                function() {
                    if ((currentColumn < indexValues.length+1)) {
                        if ([0, 3, playerIn].indexOf(gameBoardIn[currentRow][currentColumn+1]) >= 0) {
                            return true;
                        }
                    }
                    return false;
                },
                function() {
                    if (currentRow < indexValues.length+1) { 
                        if ([0, 3, playerIn].indexOf(gameBoardIn[currentRow+1][currentColumn]) >= 0) {
                            return true;
                        }
                    }
                    return false;
                },
                function() {
                    if ((currentRow < indexValues.length+1) && (currentColumn > 0)) {
                        if ([0, 3, playerIn].indexOf(gameBoardIn[currentRow+1][currentColumn-1]) >= 0) {
                            return true;
                        }
                    }
                    return false;
                },
                function() {
                    if (currentColumn > 0) {
                        if ([0, 3, playerIn].indexOf(gameBoardIn[currentRow][currentColumn-1]) >= 0) {
                            return true;
                        }
                    }
                    return false;
                },
                function() {
                    if (currentRow > 0) {
                        if ([0, 3, playerIn].indexOf(gameBoardIn[currentRow-1][currentColumn]) >= 0) {
                            return true;
                        }
                    }
                    return false;
                },
                function() {
                    if ((currentRow > 0) && (currentColumn < indexValues.length+1)) {
                        if ([0, 3, playerIn].indexOf(gameBoardIn[currentRow-1][currentColumn+1]) >= 0) {
                            return true;
                        }
                    }
                    return false;
                }
            ];

            let incrementFunctions = [
                function() {
                    currentColumn++;
                },
                function() {
                    currentRow++;
                },
                function() {
                    currentRow++;
                    currentColumn--;
                },
                function() {
                    currentColumn--;
                },
                function() {
                    currentRow--;
                },
                function() {
                    currentRow--;
                    currentColumn++;
                }
            ];

            let lastIndex = null;
            let startIndex = null;

            let currentRow = null;
            let currentColumn = null;

            if (playerIn == 1) {
                startIndex = 0;
                currentRow = 0;
                currentColumn = 0;
            } else if (playerIn == 2) {
                startIndex = 0;
                currentRow = 0;
                currentColumn = 0;
            }

            if (playerIn == 1) {
                while (true) {
                    if (!traverseFunctions[startIndex]() && traverseFunctions[((startIndex+1) %% 6)]()) {
                        startIndex = ((startIndex+1) %% 6);
                        incrementFunctions[startIndex]();
                        startIndex = ((startIndex+4) %% 6);
                    } else if (traverseFunctions[startIndex]() && !traverseFunctions[((startIndex+1) %% 6)]()) {
                        startIndex = ((startIndex+5) %% 6);
                    } else if (traverseFunctions[startIndex]() && traverseFunctions[((startIndex+1) %% 6)]()) {
                        startIndex = ((startIndex+5) %% 6);
                    } else if (!traverseFunctions[startIndex]() && !traverseFunctions[((startIndex+1) %% 6)]()) {
                        startIndex = ((startIndex+1) %% 6);
                    } else {
                        startIndex = ((startIndex+1) %% 6);
                    }
                    if (gameBoardIn[currentRow][currentColumn] == 0) {
                        break;
                    }
                }
            } else if (playerIn == 2) {
                while (true) {
                    if (traverseFunctions[startIndex]() && !traverseFunctions[((startIndex+1) %% 6)]()) {
                        incrementFunctions[startIndex]();
                    } else if (!traverseFunctions[startIndex]() && traverseFunctions[((startIndex+1) %% 6)]()) {
                        startIndex = ((startIndex+1) %% 6);
                    } else if (traverseFunctions[startIndex]() && traverseFunctions[((startIndex+1) %% 6)]()) {
                        startIndex = ((startIndex+1) %% 6);
                    } else if (!traverseFunctions[startIndex]() && !traverseFunctions[((startIndex+1) %% 6)]()) {
                        startIndex = ((startIndex+5) %% 6);
                    } else {
                        startIndex = ((startIndex+5) %% 6);
                    }
                    if (gameBoardIn[currentRow][currentColumn] == 0) {
                        break;
                    }
                }
            }

            if ((playerIn == 1) && (currentColumn >= indexValues.length+1)) {
                return true;
            } else if ((playerIn == 2) && (currentRow >= indexValues.length+1)) {
                return true;
            }

            return false;
        }

        function sendPlayerMove() {
            showMoves();
        }
        
        function drawHexPiece(row, column, color, onclick) {
            drawSVGPiece(row, column, color, onclick);
        }
        
        function drawHexCell(row, column, onclick) {
            drawSVGCell(row, column, onclick);
        }

        function getPlayerRowColumn(playerFlipCardIn, playerPlayCardIn) {

            let playerMove = {
                "row" : null,
                "column" : null
            }

            if (playerFlipCardIn) {
                if ("br".indexOf(playerFlipCardIn.substr(0,1).toLowerCase()) == 1) {
                    playerMove["row"] = indexValues.indexOf(playerFlipCardIn.substr(1))+1;
                } else {
                    playerMove["column"] = indexValues.indexOf(playerFlipCardIn.substr(1))+1;
                }
            }

            if (playerPlayCardIn) {
                if ("br".indexOf(playerPlayCardIn.substr(0,1).toLowerCase()) == 1) {
                    playerMove["row"] = indexValues.indexOf(playerPlayCardIn.substr(1))+1;
                } else {
                    playerMove["column"] = indexValues.indexOf(playerPlayCardIn.substr(1))+1;
                }
            }

            return playerMove;
        }
        
        function displayMove() {

            let playerMove = getPlayerRowColumn(playerFlipCard, playerPlayCard);
                
            if (playerFlipCard != null && playerPlayCard == null) {
                    
                if (playerMove["row"] != null) {
                    deleteLine = "r"+playerMove["row"];
                    for (let lineIndex = 1; lineIndex <= indexValues.length; lineIndex++) {
                        if (gameBoard[playerMove["row"]][lineIndex] == null) {
                            drawHexCell(playerMove["row"], lineIndex, null);
                            drawHexPiece(playerMove["row"], lineIndex, "cyan", "makeCellMove("+playerMove["row"]+","+lineIndex+",player);");
                        }
                    }
                } else if (playerMove["column"] != null) {
                    deleteLine = "c"+playerMove["column"];
                    for (let lineIndex = 1; lineIndex <= indexValues.length; lineIndex++) {
                        if (gameBoard[lineIndex][playerMove["column"]] == null) {
                            drawHexCell(lineIndex, playerMove["column"], null);
                            drawHexPiece(lineIndex, playerMove["column"], "cyan", "makeCellMove("+lineIndex+","+playerMove["column"]+",player);");
                        }
                    }
                }
                    
            } else if (playerFlipCard != null && playerPlayCard != null) {
                
                document.getElementById("playerFlipCardImg").src = cardset + "/" + playerSuits["b,r".split(",").indexOf(playerFlipCard.substr(0, 1))] + playerFlipCard.substr(1) + ".png";
                document.getElementById("playerPlayCardImg").src = cardset + "/" + playerSuits["b,r".split(",").indexOf(playerPlayCard.substr(0, 1))] + playerPlayCard.substr(1) + ".png";

                for (let lineIndex = 1; lineIndex <= indexValues.length; lineIndex++) {
                    if (deleteLine.substr(0,1) == "r") {
                        if (gameBoard[deleteLine.substr(1,deleteLine.length-1)][lineIndex] == null) {
                            let row = deleteLine.substr(1,deleteLine.length-1);
                            let column = lineIndex;
                            drawHexCell(row, column, "makeCellMove("+row+","+column+",player);");
                        }
                    } else if (deleteLine.substr(0,1) == "c") {
                        if (gameBoard[lineIndex][deleteLine.substr(1,deleteLine.length-1)] == null) {
                            let row = lineIndex;
                            let column = deleteLine.substr;(1,deleteLine.length-1)
                            drawHexCell(row, column, "makeCellMove("+row+","+column+",player);");
                        }
                    }
                }

                console.log("Hello?");
                drawHexCell(playerMove["row"], playerMove["column"], null);
                drawHexPiece(playerMove["row"], playerMove["column"], "cyan", "sendPlayerMove();");
            }
        }
        
        function svgHexCell(row, column, onclick) {
        
            let hexObj = document.createElementNS(svgNS, "polygon");
            let pointsStr = "";
            
            let xCoord = getXCoord(row, column);
            let yCoord = getYCoord(row, column);
            
                        pointsStr += (xCoord+(10.243311227557877))+","+(yCoord+(-5.913978494623655));
            pointsStr += " "
            pointsStr += (xCoord+(0.000000000000001))+","+(yCoord+(-11.827956989247312));
            pointsStr += " "
            pointsStr += (xCoord+(-10.243311227557875))+","+(yCoord+(-5.913978494623660));
            pointsStr += " "
            pointsStr += (xCoord+(-10.243311227557879))+","+(yCoord+(5.913978494623652));
            pointsStr += " "
            pointsStr += (xCoord+(-0.000000000000002))+","+(yCoord+(11.827956989247312));
            pointsStr += " "
            pointsStr += (xCoord+(10.243311227557879))+","+(yCoord+(5.913978494623652));

        
            hexObj.setAttribute("points", pointsStr);

            if (onclick) {
                hexObj.setAttribute("onclick", onclick);
            }
            
            hexObj.style.fill="white";
            hexObj.style.stroke="black";
        
            return hexObj;
        }
        
        function svgHexPiece(row, column, color, onclick) {
            
            let pieceObj = document.createElementNS(svgNS, "circle");
            
            pieceObj.setAttribute("cx", getXCoord(row, column));
            pieceObj.setAttribute("cy", getYCoord(row, column));
            pieceObj.setAttribute("r", 8.279569892473118);
            pieceObj.setAttribute("stroke", "black");
            pieceObj.setAttribute("fill", color);

            if (onclick) {
                pieceObj.setAttribute("onclick", onclick);
            }
            
            return pieceObj;
        }
        
        function drawSVGPiece(row, column, color, onclick) {
            let svgObj = document.getElementById(svgID);
            svgObj.appendChild(svgHexPiece(row, column, color, onclick));
        }
        
        function drawSVGCell(row, column, onclick) {
            let svgObj = document.getElementById(svgID);
            svgObj.appendChild(svgHexCell(row, column, onclick));
        }

        function startGame() {

            voltageBoardIndexes = null;

            gameBoard = [];
            voltageBoardIndexes = [];

            let voltageIndex = 0;
            let computerMove = null;

            document.getElementById("playerPlayCardImg").src = cardset + "/blank.png";
            document.getElementById("playerFlipCardImg").src = cardset + "/blank.png";
            document.getElementById("opponentPlayCardImg").src = cardset + "/blank.png";
            document.getElementById("opponentFlipCardImg").src = cardset + "/blank.png";

            for (let rowIndex = 0; rowIndex <= indexValues.length+1; rowIndex++) {
                gameBoard.push([]);
                voltageBoardIndexes.push([]);
                for (let columnIndex = 0; columnIndex <= indexValues.length+1; columnIndex++) {
                    if (
                        (rowIndex == 0 && columnIndex == indexValues.length+1) ||
                        (rowIndex == indexValues.length+1 && columnIndex == 0)
                    ) {
                        gameBoard[rowIndex].push(0);
                    } else if (
                        (rowIndex == 0 && columnIndex == 0) ||
                        (rowIndex == indexValues.length+1 && columnIndex == indexValues.length+1)
                    ) {
                        gameBoard[rowIndex].push(0);
                    } else if (rowIndex == 0 || rowIndex == indexValues.length+1) {
                        gameBoard[rowIndex].push(2);
                    } else if (columnIndex == 0 || columnIndex == indexValues.length+1) {
                        gameBoard[rowIndex].push(1);
                    } else {
                        gameBoard[rowIndex].push(null);
                        drawHexCell(rowIndex, columnIndex, "makeCellMove("+rowIndex+","+columnIndex+",player);");
                    }
                    voltageBoardIndexes[rowIndex].push(voltageIndex++);
                }
            }

            playerFlipCard = null;
            playerPlayCard = null;
            computerFlipCard = null;
            computerPlayCard = null;

            document.getElementById("loader").classList.add("loaderanimation");
            setTimeout(
                function() {
                    computerMove = calculateComputerMove(computer, gameBoard);
                    makeCellMove(computerMove["row"], computerMove["column"], computer);
                    showMoves();

                    document.getElementById("loader").classList.remove("loaderanimation");
                },
                loaderTimeout
            );
        }

        function calculateComputerMove(playerIn, gameBoardIn) {

            let maxProbability = 0;
            let maxIndex = 2;
            let returnCoords = {};

            let computerMoveSets = [
                normalizeMoves(calculateMonteCarloMoves(playerIn, gameBoardIn))[0],
                normalizeMoves(calculateVoltageMoves(playerIn, gameBoardIn))[0],
                normalizeMoves(calculateSemiBestMoves(playerIn, gameBoardIn))[0]
            ];

            for (let setIndex = 0; setIndex < computerMoveSets.length; setIndex++) {
                if (computerMoveSets[setIndex]) {
                    let currentProbability = testUsingMonteCarlo(playerIn, gameBoardIn, computerMoveSets[setIndex]["row"], computerMoveSets[setIndex]["column"]);
                    if (currentProbability > maxProbability) {
                        maxProbability = currentProbability;
                        maxIndex = setIndex;
                    }
                }
            }

            console.log(maxIndex, computerMoveSets[maxIndex]["row"], computerMoveSets[maxIndex]["column"], maxProbability);

            returnCoords = {
                "row": computerMoveSets[maxIndex]["row"],
                "column": computerMoveSets[maxIndex]["column"]
            };

            return returnCoords;
        }

        function normalizeMoves(computerMovesIn) {

            let normalizedMoves = [];
            let weights = [];

            for (let moveIndex = 0; moveIndex < computerMovesIn.length; moveIndex++) {
                weights.push(computerMovesIn[moveIndex]["weight"]);
            }

            for (let moveIndex = 0; moveIndex < computerMovesIn.length; moveIndex++) {
                let moveDict = {
                    "row": computerMovesIn[moveIndex]["row"],
                    "column": computerMovesIn[moveIndex]["column"],
                    "zscore": ((computerMovesIn[moveIndex]["weight"] - math.mean(weights)) / math.std(weights))
                }
                normalizedMoves.push(moveDict);
            }

            normalizedMoves = normalizedMoves.sort(
                function(a, b) {
                    return b["zscore"] -  a["zscore"];
                }
            );

            return normalizedMoves;
        }

        function calculateSemiBestMoves(playerIn, gameBoardIn) {

            let computerMoves = [];

            let voltageBoardVectors = [null, calculateVoltages(1, gameBoardIn), calculateVoltages(2, gameBoardIn)];

            if (voltageBoardVectors[1].length && voltageBoardVectors[2].length) {
                computerMoves = getSemiBestMoves(playerIn, gameBoardIn, calculateVoltageDeltas(voltageBoardVectors[1]), calculateVoltageDeltas(voltageBoardVectors[2]));
            }

            return computerMoves;
        }

        function testUsingMonteCarlo(playerIn, gameBoardIn, rowIn, columnIn) {

            let copyBoard = [];
            let calcIndex1 = 0;

            let winCount = 0;

            let opponent = (playerIn %% 2) + 1;

            for (let copyIndex = 0; copyIndex < gameBoardIn.length; copyIndex++) {
                copyBoard.push(Object.assign([], gameBoardIn[copyIndex]));
            }

            copyBoard[rowIn][columnIn] = playerIn;
            let emptyCoords = getEmptyCoords(copyBoard);

            while (calcIndex1 < randomMoveLimit) {

                let randomBoard = randomizeBoardMoves(copyBoard, emptyCoords, playerIn);

                let test1 = null;
                let test2 = null;

                if (evaluateWin(playerIn, randomBoard)) {
                    test1 = playerIn;
                } else if (evaluateWin(opponent, randomBoard)) {
                    test1 = opponent;
                }

                randomBoard[rowIn][columnIn] = opponent;

                if (evaluateWin(playerIn, randomBoard)) {
                    test2 = playerIn;
                } else if (evaluateWin(opponent, randomBoard)) {
                    test2 = opponent;
                }

                if (test1 != test2) {
                    winCount++;
                }

                calcIndex1++;
            }
            
            return (winCount / randomMoveLimit);
        }

        function calculateMonteCarloMoves(playerIn, gameBoardIn) {

            let computerMoves = [];

            let emptyCoords = getEmptyCoords(gameBoardIn);
            let coordCounts = [];

            let opponent = (playerIn %% 2) + 1;

            for (let initIndex = 0; initIndex < emptyCoords.length; initIndex++) {
                coordCounts.push(0);
            }

            let calcIndex1 = 0;

            while (calcIndex1 < randomMoveLimit) {

                let randomBoard = randomizeBoardMoves(gameBoardIn, emptyCoords);

                for (let calcIndex2 = 0; calcIndex2 < emptyCoords.length; calcIndex2++) {

                    let test1 = null;
                    let test2 = null;

                    if (evaluateWin(playerIn, randomBoard)) {
                        test1 = playerIn;
                    } else if (evaluateWin(opponent, randomBoard)) {
                        test1 = opponent;
                    }

                    randomBoard[emptyCoords[calcIndex2][0]][emptyCoords[calcIndex2][1]] = (randomBoard[emptyCoords[calcIndex2][0]][emptyCoords[calcIndex2][1]] %% 2) + 1;

                    if (evaluateWin(playerIn, randomBoard)) {
                        test2 = playerIn;
                    } else if (evaluateWin(opponent, randomBoard)) {
                        test2 = opponent;
                    }

                    if (test1 != test2) {
                        coordCounts[calcIndex2]++;
                    }
                }

                calcIndex1++;
            }

            for (let checkIndex = 0; checkIndex < emptyCoords.length; checkIndex++) {
                let computerMove = {
                    "row": emptyCoords[checkIndex][0],
                    "column": emptyCoords[checkIndex][1],
                    "weight": coordCounts[checkIndex]
                }
                computerMoves.push(computerMove);
            }

            return computerMoves;
        }

        function getEmptyCoords(gameBoardIn) {

            let emptyCoords = [];

            for (let rowIndex = 1; rowIndex <= indexValues.length; rowIndex++) {
                for (let columnIndex = 1; columnIndex <= indexValues.length; columnIndex++) {
                    if (gameBoardIn[rowIndex][columnIndex] == null) {
                        emptyCoords.push([rowIndex, columnIndex]);
                    }
                }
            }

            return emptyCoords;
        }

        function randomizeBoardMoves(gameBoardIn, emptyCoordsIn, playerIn) {

            let randomBoard = [];
            let randomPlayer = (Math.random() < 0.5, 1, 2);
            let populateCoords = [];

            let opponent = (playerIn %% 2) + 1;

            if (playerIn) {
                randomPlayer = opponent;
            }

            for (let copyIndex = 0; copyIndex < gameBoardIn.length; copyIndex++) {
                randomBoard.push(Object.assign([], gameBoardIn[copyIndex]));
            }

            populateCoords = Object.assign([], emptyCoordsIn);
            for (let moveIndex = populateCoords.length; moveIndex > 0; moveIndex--) {
                let moveChoice = Math.floor(Math.random() * moveIndex);
                randomBoard[populateCoords[moveChoice][0]][populateCoords[moveChoice][1]] = randomPlayer;
                populateCoords.splice(moveChoice, 1);
                randomPlayer = (randomPlayer %% 2) + 1;
            }

            return randomBoard;
        }

        function calculateVoltageMoves(playerIn, gameBoardIn) {

            let computerMoves = [];

            let voltageBoardVectors = [null, calculateVoltages(1, gameBoardIn), calculateVoltages(2, gameBoardIn)];

            if (voltageBoardVectors[1].length && voltageBoardVectors[2].length) {
                computerMoves = evaluateMoves(playerIn, gameBoardIn, calculateVoltageDeltas(voltageBoardVectors[1]), calculateVoltageDeltas(voltageBoardVectors[2]));
            }

            return computerMoves;
        }

        function resetVoltageSolutions(gameBoardIn) {

            let voltageBoardSolutions = [];

            for (let rowIndex = 0; rowIndex < gameBoardIn.length; rowIndex++) {
                for (let columnIndex = 0; columnIndex < gameBoardIn[rowIndex].length; columnIndex++) {
                    voltageBoardSolutions.push(0);
                }
            }

            return voltageBoardSolutions;
        }

        function resetVoltageBoards(solutionsLengthIn) {

            let voltageBoard = [];

            for (let rowIndex = 0; rowIndex < solutionsLengthIn; rowIndex++) {
                voltageBoard.push([]);
                for (let colIndex = 0; colIndex <  solutionsLengthIn; colIndex++) {
                    voltageBoard[rowIndex].push(0);
                }
            }

            return voltageBoard;
        }

        function calculateVoltages(playerIn, gameBoardIn) {

            let returnVoltages = [];

            let voltageBoardSolutions = resetVoltageSolutions(gameBoardIn);
            let voltageBoards = resetVoltageBoards(voltageBoardSolutions.length);

            for (let rowIndex = 0; rowIndex < gameBoardIn.length; rowIndex++) {
                for (let columnIndex = 0; columnIndex < gameBoardIn[rowIndex].length; columnIndex++) {

                    if (rowIndex == 0 && columnIndex == 0) {
                        continue;
                    } else if (rowIndex == indexValues.length+1 && columnIndex == indexValues.length+1) {
                        continue;
                    } else if (rowIndex == 0 && playerIn == 2) {
                        voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex]] = 1;
                        voltageBoardSolutions[voltageBoardIndexes[rowIndex][columnIndex]] = testVoltage;
                    } else if (columnIndex == 0 && playerIn == 1) {
                        voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex]] = 1;
                        voltageBoardSolutions[voltageBoardIndexes[rowIndex][columnIndex]] = testVoltage;
                    } else if (rowIndex == 0 && playerIn == 1) {
                        voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex]] = 1;
                        voltageBoardSolutions[voltageBoardIndexes[rowIndex][columnIndex]] = 0;
                    } else if (columnIndex == 0 && playerIn == 2) {
                        voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex]] = 1;
                        voltageBoardSolutions[voltageBoardIndexes[rowIndex][columnIndex]] = 0;
                    } else if (rowIndex == gameBoardIn.length-1 || columnIndex == gameBoardIn[rowIndex].length-1) {
                        voltageBoards[playerIn][voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex]] = 1;
                        voltageBoardSolutions[voltageBoardIndexes[rowIndex][columnIndex]] = 0;
                    } else {

                        let connCount = 0;
                        let sameCount = 0;
                        let signCount = -1;

                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex-1, columnIndex, gameBoardIn, playerIn) == 2) {
                            sameCount++;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex-1, columnIndex+1, gameBoardIn, playerIn) == 2) {
                            sameCount++;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex, columnIndex-1, gameBoardIn, playerIn) == 2) {
                            sameCount++;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex, columnIndex+1, gameBoardIn, playerIn) == 2) {
                            sameCount++;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex+1, columnIndex-1, gameBoardIn, playerIn) == 2) {
                            sameCount++;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex+1, columnIndex, gameBoardIn, playerIn) == 2) {
                            sameCount++;
                        }

                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex-1, columnIndex, gameBoardIn, playerIn) == 2) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex-1][columnIndex]] = -2/sameCount;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex-1, columnIndex+1, gameBoardIn, playerIn) == 2) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex-1][columnIndex+1]] -2/sameCount;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex, columnIndex-1, gameBoardIn, playerIn) == 2) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex-1]] = -2/sameCount;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex, columnIndex+1, gameBoardIn, playerIn) == 2) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex+1]] = -2/sameCount;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex+1, columnIndex-1, gameBoardIn, playerIn) == 2) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex+1][columnIndex-1]] = -2/sameCount;
                        }
                        if (getVoltageRelationship(rowIndex, columnIndex, rowIndex+1, columnIndex, gameBoardIn, playerIn) == 2) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex+1][columnIndex]] = -2/sameCount;
                        }

                        if (sameCount) {
                            signCount = 1;
                        }

                        if ([1, 3].indexOf(getVoltageRelationship(rowIndex, columnIndex, rowIndex-1, columnIndex, gameBoardIn, playerIn)) > -1) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex-1][columnIndex]] = signCount;
                            connCount++;
                        }
                        if ([1, 3].indexOf(getVoltageRelationship(rowIndex, columnIndex, rowIndex-1, columnIndex+1, gameBoardIn, playerIn)) > -1) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex-1][columnIndex+1]] = signCount;
                            connCount++;
                        }
                        if ([1, 3].indexOf(getVoltageRelationship(rowIndex, columnIndex, rowIndex, columnIndex-1, gameBoardIn, playerIn)) > -1) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex-1]] = signCount;
                            connCount++;
                        }
                        if ([1, 3].indexOf(getVoltageRelationship(rowIndex, columnIndex, rowIndex, columnIndex+1, gameBoardIn, playerIn)) > -1) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex+1]] = signCount;
                            connCount++;
                        }
                        if ([1, 3].indexOf(getVoltageRelationship(rowIndex, columnIndex, rowIndex+1, columnIndex-1, gameBoardIn, playerIn)) > -1) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex+1][columnIndex-1]] = signCount;
                            connCount++;
                        }
                        if ([1, 3].indexOf(getVoltageRelationship(rowIndex, columnIndex, rowIndex+1, columnIndex, gameBoardIn, playerIn)) > -1) {
                            voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex+1][columnIndex]] = signCount;
                            connCount++;
                        }

                        voltageBoards[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex]] = connCount;
                        voltageBoardSolutions[voltageBoardIndexes[rowIndex][columnIndex]] = 0;
                    }
                }
            }

            try {
                let sparseVoltageMatrix = numeric.ccsSparse(voltageBoards);
                returnVoltages = numeric.ccsLUPSolve(numeric.ccsLUP(sparseVoltageMatrix), voltageBoardSolutions);
            } catch (err) {
                console.log(err);
            }
            
            return returnVoltages;
        }

        function getVoltageRelationship(refRow, refCol, calcRow, calcCol, gameBoardIn, playerIn) {

            if (gameBoardIn[refRow][refCol] == null && gameBoardIn[calcRow][calcCol] == playerIn) {
                return 1;
            } else if (gameBoardIn[refRow][refCol] == playerIn && gameBoardIn[calcRow][calcCol] == null) {
                return 1;
            } else if (gameBoardIn[refRow][refCol] == null && gameBoardIn[calcRow][calcCol] == 0) {
                return 1;
            } else if (gameBoardIn[refRow][refCol] == 0 && gameBoardIn[calcRow][calcCol] == null) {
                return 1;
            } else if (gameBoardIn[refRow][refCol] == playerIn && gameBoardIn[calcRow][calcCol] == playerIn) {
                return 2;
            } else if (gameBoardIn[refRow][refCol] == 0 && gameBoardIn[calcRow][calcCol] == playerIn) {
                return 2;
            } else if (gameBoardIn[refRow][refCol] == playerIn && gameBoardIn[calcRow][calcCol] == 0) {
                return 2;
            } else if (gameBoardIn[refRow][refCol] == null && gameBoardIn[calcRow][calcCol] == null) {
                return 3;
            }

            return 0;
        }

        function calculateVoltageDeltas(voltageBoardVectors) {

            let deltasOut = [];

            for (let rowIndex = 0; rowIndex < voltageBoardVectors.length; rowIndex++) {
                deltasOut.push([]);
                for (let columnIndex = 0; columnIndex < voltageBoardVectors.length; columnIndex++) {
                    deltasOut[rowIndex].push(0);
                }
            }

            for (let rowIndex = 1; rowIndex <= indexValues.length; rowIndex++) {
                for (let columnIndex = 1; columnIndex <= indexValues.length; columnIndex++) {
                    deltasOut[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex-1][columnIndex]] = voltageBoardVectors[voltageBoardIndexes[rowIndex][columnIndex]] - voltageBoardVectors[voltageBoardIndexes[rowIndex-1][columnIndex]];
                    deltasOut[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex-1][columnIndex+1]] = voltageBoardVectors[voltageBoardIndexes[rowIndex][columnIndex]] - voltageBoardVectors[voltageBoardIndexes[rowIndex-1][columnIndex+1]];
                    deltasOut[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex-1]] = voltageBoardVectors[voltageBoardIndexes[rowIndex][columnIndex]] - voltageBoardVectors[voltageBoardIndexes[rowIndex][columnIndex-1]];
                    deltasOut[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex][columnIndex+1]] = voltageBoardVectors[voltageBoardIndexes[rowIndex][columnIndex]] - voltageBoardVectors[voltageBoardIndexes[rowIndex][columnIndex+1]];
                    deltasOut[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex+1][columnIndex-1]] = voltageBoardVectors[voltageBoardIndexes[rowIndex][columnIndex]] - voltageBoardVectors[voltageBoardIndexes[rowIndex+1][columnIndex-1]];
                    deltasOut[voltageBoardIndexes[rowIndex][columnIndex]][voltageBoardIndexes[rowIndex+1][columnIndex]] = voltageBoardVectors[voltageBoardIndexes[rowIndex][columnIndex]] - voltageBoardVectors[voltageBoardIndexes[rowIndex+1][columnIndex]];
                }
            }

            return deltasOut;
        }

        function evaluateMoveDeltas(rowIndexIn, columnIndexIn, playerIn, deltasIn) {

            let moveEval = 0;

            if ([1, 3].indexOf(getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn, gameBoard, playerIn)) > -1) {
                moveEval += deltasIn[voltageBoardIndexes[rowIndexIn][columnIndexIn]][voltageBoardIndexes[rowIndexIn-1][columnIndexIn]];
            } else if ([1, 3].indexOf(getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+1, gameBoard, playerIn)) > -1) {
                moveEval += deltasIn[voltageBoardIndexes[rowIndexIn][columnIndexIn]][voltageBoardIndexes[rowIndexIn-1][columnIndexIn+1]];
            } else if ([1, 3].indexOf(getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn-1, gameBoard, playerIn)) > -1) {
                moveEval += deltasIn[voltageBoardIndexes[rowIndexIn][columnIndexIn]][voltageBoardIndexes[rowIndexIn][columnIndexIn-1]];
            } else if ([1, 3].indexOf(getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn+1, gameBoard, playerIn)) > -1) {
                moveEval += deltasIn[voltageBoardIndexes[rowIndexIn][columnIndexIn]][voltageBoardIndexes[rowIndexIn][columnIndexIn+1]];
            } else if ([1, 3].indexOf(getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-1, gameBoard, playerIn)) > -1) {
                moveEval += deltasIn[voltageBoardIndexes[rowIndexIn][columnIndexIn]][voltageBoardIndexes[rowIndexIn+1][columnIndexIn-1]];
            } else if ([1, 3].indexOf(getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn, gameBoard, playerIn)) > -1) {
                moveEval += deltasIn[voltageBoardIndexes[rowIndexIn][columnIndexIn]][voltageBoardIndexes[rowIndexIn+1][columnIndexIn]];
            }

            return moveEval;
        }

        function moveBordersOpponent(rowIndexIn, columnIndexIn, gameStateIn, playerIn) {

            let moveEval = false;
            let opponent = (playerIn %% 2) + 1;

            if (
                0 >= rowIndexIn-1 ||
                gameStateIn.length <= rowIndexIn+1 ||
                0 >= columnIndexIn-1 ||
                gameStateIn.length <= columnIndexIn+1
            ) {
                moveEval = false;
            } else if (!getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn, gameStateIn, playerIn) && getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn -1, columnIndexIn, gameStateIn, opponent)) {
                moveEval = true;
            } else if (!getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+1, gameStateIn, playerIn) && getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+1, gameStateIn, opponent)) {
                moveEval = true;
            } else if (!getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn-1, gameStateIn, playerIn) && getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn-1, gameStateIn, opponent)) {
                moveEval = true;
            } else if (!getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn+1, gameStateIn, playerIn) && getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn+1, gameStateIn, opponent)) {
                moveEval = true;
            } else if (!getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-1, gameStateIn, playerIn) && getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-1, gameStateIn, opponent)) {
                moveEval = true;
            } else if (!getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn, gameStateIn, playerIn) && getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn, gameStateIn, opponent)) {
                moveEval = true;
            }

            return moveEval;
        }

        function moveFitsOpponentDiamond(rowIndexIn, columnIndexIn, gameStateIn, playerIn) {

            let moveEval = false;
            let opponent = (playerIn %% 2) + 1;

            if (
                0 >= rowIndexIn-2 ||
                gameStateIn.length <= rowIndexIn+2 ||
                0 >= columnIndexIn-2 ||
                gameStateIn.length <= columnIndexIn+2
            ) {
                moveEval = false;
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn+1, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn+1, gameStateIn, opponent) == 1
            ) {
                moveEval = true;
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn-1, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn-1, gameStateIn, opponent) == 1
            ) {
                moveEval = true;
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+2, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn+1, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+1, gameStateIn, opponent) == 1
            ) {
                moveEval = true;
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-2, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn-1, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-1, gameStateIn, opponent) == 1
            ) {
                moveEval = true;
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+2, columnIndexIn-1, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-1, gameStateIn, opponent) == 1
            ) {
                moveEval = true;
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-2, columnIndexIn+1, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn, gameStateIn, opponent) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+1, gameStateIn, opponent) == 1
            ) {
                moveEval = true;
            }

            return moveEval;
        }

        function moveFitsEstablishedDiamond(rowIndexIn, columnIndexIn, gameStateIn, playerIn) {

            let moveEval = [];

            if (
                0 >= rowIndexIn-2 ||
                gameStateIn.length <= rowIndexIn+2 ||
                0 >= columnIndexIn-2 ||
                gameStateIn.length <= columnIndexIn+2
            ) {
                moveEval = [];
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn+1, gameStateIn, playerIn) == 2 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn, gameStateIn, playerIn) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn+1, gameStateIn, playerIn) == 1
            ) {
                moveEval.push({"row" : rowIndexIn+1, "column" : columnIndexIn});
                moveEval.push({"row" : rowIndexIn, "column" : columnIndexIn+1});
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn-1, gameStateIn, playerIn) == 2 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn, gameStateIn, playerIn) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn-1, gameStateIn, playerIn) == 1) {
                moveEval.push({"row" : rowIndexIn-1, "column" : columnIndexIn});
                moveEval.push({"row" : rowIndexIn, "column" : columnIndexIn-1});
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+2, gameStateIn, playerIn) == 2 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn+1, gameStateIn, playerIn) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+1, gameStateIn, playerIn) == 1
            ) {
                moveEval.push({"row" : rowIndexIn, "column" : columnIndexIn+1});
                moveEval.push({"row" : rowIndexIn-1, "column" : columnIndexIn+1});
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-2, gameStateIn, playerIn) == 2 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn, columnIndexIn-1, gameStateIn, playerIn) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-1, gameStateIn, playerIn) == 1
            ) {
                moveEval.push({"row" : rowIndexIn, "column" : columnIndexIn-1});
                moveEval.push({"row" : rowIndexIn+1, "column" : columnIndexIn-1});
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+2, columnIndexIn-1, gameStateIn, playerIn) == 2 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn, gameStateIn, playerIn) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn+1, columnIndexIn-1, gameStateIn, playerIn) == 1
            ) {
                moveEval.push({"row" : rowIndexIn+1, "column" : columnIndexIn});
                moveEval.push({"row" : rowIndexIn+1, "column" : columnIndexIn-1});
            } else if (
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-2, columnIndexIn+1, gameStateIn, playerIn) == 2 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn, gameStateIn, playerIn) == 1 &&
                getVoltageRelationship(rowIndexIn, columnIndexIn, rowIndexIn-1, columnIndexIn+1, gameStateIn, playerIn) == 1
            ) {
                moveEval.push({"row" : rowIndexIn-1, "column" : columnIndexIn});
                moveEval.push({"row" : rowIndexIn-1, "column" : columnIndexIn+1});
            }

            return moveEval;
        }

        function moveAdvantageousOpponent(rowIndexIn, columnIndexIn, gameStateIn, playerIn) {

            let moveEval = false;

            if (moveFitsOpponentDiamond(rowIndexIn, columnIndexIn, gameStateIn, playerIn)) {
                moveEval = true;
            } else if (moveBordersOpponent(rowIndexIn, columnIndexIn, gameStateIn, playerIn)) {
                moveEval = true;
            }

            return moveEval;
        }

        function getDefensiveMoves(playerIn, gameStateIn, playerFactorsIn, opponentFactorsIn) {

            let opponent = (playerIn %% 2) + 1;
            
            let playerMovesReturn = [];

            if (!playerMovesReturn.length) {
                for (let rowIndex = 0; rowIndex < gameStateIn.length; rowIndex++) {
                    for (let columnIndex = 0; columnIndex < gameStateIn[rowIndex].length; columnIndex++) {
                        if (gameStateIn[rowIndex][columnIndex] == opponent) {
                            let diamonds = moveFitsEstablishedDiamond(rowIndex, columnIndex, gameStateIn, opponent);
                            if (diamonds.length) {
                                for (let moveIndex = 0; moveIndex < diamonds.length; moveIndex+=2) {
                                    if (gameStateIn[diamonds[moveIndex]["row"]][diamonds[moveIndex]["column"]] == null && gameStateIn[diamonds[moveIndex+1]["row"]][diamonds[moveIndex+1]["column"]] == null) {
                                        playerMovesReturn.push({
                                            "playerFactor" : playerFactorsIn[diamonds[moveIndex]["row"]][diamonds[moveIndex]["column"]],
                                            "opponentFactor" : opponentFactorsIn[diamonds[moveIndex]["row"]][diamonds[moveIndex]["column"]],
                                            "normPlayerFactor" : null,
                                            "normOpponentFactor" : null,
                                            "weight" : null,
                                            "row" : diamonds[moveIndex]["row"]+1,
                                            "column" : diamonds[moveIndex]["column"]+1
                                        });
                                        playerMovesReturn.push({
                                            "playerFactor" : playerFactorsIn[diamonds[moveIndex+1]["row"]][diamonds[moveIndex+1]["column"]],
                                            "opponentFactor" : opponentFactorsIn[diamonds[moveIndex+1]["row"]][diamonds[moveIndex+1]["column"]],
                                            "normPlayerFactor" : null,
                                            "normOpponentFactor" : null,
                                            "weight" : null,
                                            "row" : diamonds[moveIndex+1]["row"]+1,
                                            "column" : diamonds[moveIndex+1]["column"]+1
                                        });
                                    }
                                }
                            }
                        }
                    }
                }
            }

            for (let rowIndex = 0; rowIndex < gameStateIn.length; rowIndex++) {
                for (let columnIndex = 0; columnIndex < gameStateIn[rowIndex].length; columnIndex++) {
                    if (gameStateIn[rowIndex][columnIndex] == null) {
                        if (moveAdvantageousOpponent(rowIndex, columnIndex, gameStateIn, playerIn)) {
                            playerMovesReturn.push({
                                "playerFactor" : playerFactorsIn[rowIndex][columnIndex],
                                "opponentFactor" : opponentFactorsIn[rowIndex][columnIndex],
                                "normPlayerFactor" : null,
                                "normOpponentFactor" : null,
                                "weight" : null,
                                "row" : rowIndex+1,
                                "column" : columnIndex+1
                            });
                        }
                    }
                }
            }

            return playerMovesReturn;
        }

        function setupMovesReturn(moveReturnIn) {

            let moveReturnOut = moveReturnIn;

            let playerMoveCount = 0;
            let playerMovesSum = 0;
            let playerMaxFactor = 0;

            let opponentMoveCount = 0;
            let opponentMovesSum = 0;
            let opponentMaxFactor = 0;

            let ratioNormCount = 0;
            let ratioNormsSum = 0;

            for (let moveIndex = 0; moveIndex < moveReturnOut.length; moveIndex++) {
                if (moveReturnOut[moveIndex]["playerFactor"] != null) {
                    playerMoveCount++;
                    playerMovesSum += moveReturnOut[moveIndex]["playerFactor"];
                    playerMaxFactor = (playerMaxFactor < moveReturnOut[moveIndex]["playerFactor"] ? moveReturnOut[moveIndex]["playerFactor"] : playerMaxFactor);
                }
                if (moveReturnOut[moveIndex]["opponentFactor"] != null) {
                    opponentMoveCount++;
                    opponentMovesSum += moveReturnOut[moveIndex]["opponentFactor"];
                    opponentMaxFactor = (opponentMaxFactor < moveReturnOut[moveIndex]["opponentFactor"] ? moveReturnOut[moveIndex]["opponentFactor"] : opponentMaxFactor);
                }
            }

            for (let moveIndex = 0; moveIndex < moveReturnOut.length; moveIndex++) {

                if (moveReturnOut[moveIndex]["playerFactor"] != null && moveReturnOut[moveIndex]["opponentFactor"] != null) {
                    moveReturnOut[moveIndex]["normPlayerFactor"] = moveReturnOut[moveIndex]["playerFactor"]*playerMoveCount/playerMovesSum;
                    moveReturnOut[moveIndex]["normOpponentFactor"] = moveReturnOut[moveIndex]["opponentFactor"]*opponentMoveCount/opponentMovesSum;
                    moveReturnOut[moveIndex]["weight"] = moveReturnOut[moveIndex]["normPlayerFactor"]/moveReturnOut[moveIndex]["normOpponentFactor"];
                    ratioNormCount++;
                    ratioNormsSum += moveReturnOut[moveIndex]["weight"];
                }

                if (moveReturnOut[moveIndex]["playerFactor"] == null) {
                    moveReturnOut[moveIndex]["playerFactor"] = playerMaxFactor;
                }
                if (moveReturnOut[moveIndex]["opponentFactor"] == null) {
                    moveReturnOut[moveIndex]["opponentFactor"] = opponentMaxFactor;
                }
                if (moveReturnOut[moveIndex]["weight"] == null) {
                    moveReturnOut[moveIndex]["weight"] = 0;
                }
            }

            moveReturnOut.sort(
                function(a, b) {
                    return a["weight"] - b["weight"];
                }
            );

            return moveReturnOut;
        }

        function evaluateMoves(playerIn, gameBoardIn, playerDeltas, opponentDeltas) {

            let opponent = (playerIn %% 2) + 1;

            let gameState = [];

            let playerFactors = [];
            let opponentFactors = [];

            let playerMoves = [];
            let playerMoveReturn = null;
            let playerMovesOut = [];

            let playerEvals = [];
            let opponentEvals = [];

            for (let rowIndex = 1; rowIndex <= indexValues.length; rowIndex++) {
                playerFactors.push([]);
                opponentFactors.push([]);
                gameState.push([]);
                for (let columnIndex = 1; columnIndex <= indexValues.length; columnIndex++) {
                    playerFactors[rowIndex-1].push(null);
                    opponentFactors[rowIndex-1].push(null);
                    gameState[rowIndex-1].push(null);
                    if (gameBoardIn[rowIndex][columnIndex] == null) {

                        playerEvals.push(evaluateMoveDeltas(rowIndex, columnIndex, playerIn, playerDeltas));
                        opponentEvals.push(evaluateMoveDeltas(rowIndex, columnIndex, opponent, opponentDeltas));

                        playerEvals[playerEvals.length-1] *= playerEvals[playerEvals.length-1];
                        playerEvals[playerEvals.length-1] /= 6;
                        playerEvals[playerEvals.length-1] = Math.sqrt(playerEvals[playerEvals.length-1]);

                        opponentEvals[opponentEvals.length-1] *= opponentEvals[opponentEvals.length-1];
                        opponentEvals[opponentEvals.length-1] /= 6;
                        opponentEvals[opponentEvals.length-1] = Math.sqrt(opponentEvals[opponentEvals.length-1]);

                        playerFactors[rowIndex-1][columnIndex-1] = (playerEvals[playerEvals.length-1] ? opponentEvals[opponentEvals.length-1]/playerEvals[playerEvals.length-1] : null);
                        opponentFactors[rowIndex-1][columnIndex-1] = (opponentEvals[opponentEvals.length-1] ? playerEvals[playerEvals.length-1]/opponentEvals[opponentEvals.length-1] : null);

                        playerMoves.push({
                            "playerFactor" : playerFactors[rowIndex-1][columnIndex-1],
                            "opponentFactor" : opponentFactors[rowIndex-1][columnIndex-1],
                            "normPlayerFactor" : null,
                            "normOpponentFactor" : null,
                            "weight" : null,
                            "row" : rowIndex,
                            "column" : columnIndex
                        });

                    } else {
                        gameState[rowIndex-1][columnIndex-1] = gameBoardIn[rowIndex][columnIndex];
                    }
                }
            }

            playerMovesOut = setupMovesReturn(playerMoves);

            return playerMovesOut;
        }

        function getSemiBestMoves(playerIn, gameBoardIn, playerDeltas, opponentDeltas) {

            let opponent = (playerIn %% 2) + 1;

            let gameState = [];

            let playerFactors = [];
            let opponentFactors = [];

            let playerMoves = [];
            let playerMoveReturn = null;
            let playerMovesOut = [];

            let playerEvals = [];
            let opponentEvals = [];

            for (let rowIndex = 1; rowIndex <= indexValues.length; rowIndex++) {
                playerFactors.push([]);
                opponentFactors.push([]);
                gameState.push([]);
                for (let columnIndex = 1; columnIndex <= indexValues.length; columnIndex++) {
                    playerFactors[rowIndex-1].push(null);
                    opponentFactors[rowIndex-1].push(null);
                    gameState[rowIndex-1].push(null);
                    if (gameBoardIn[rowIndex][columnIndex] == null) {

                        playerEvals.push(evaluateMoveDeltas(rowIndex, columnIndex, playerIn, playerDeltas));
                        opponentEvals.push(evaluateMoveDeltas(rowIndex, columnIndex, opponent, opponentDeltas));

                        playerEvals[playerEvals.length-1] *= playerEvals[playerEvals.length-1];
                        playerEvals[playerEvals.length-1] /= 6;
                        playerEvals[playerEvals.length-1] = Math.sqrt(playerEvals[playerEvals.length-1]);

                        opponentEvals[opponentEvals.length-1] *= opponentEvals[opponentEvals.length-1];
                        opponentEvals[opponentEvals.length-1] /= 6;
                        opponentEvals[opponentEvals.length-1] = Math.sqrt(opponentEvals[opponentEvals.length-1]);

                        playerFactors[rowIndex-1][columnIndex-1] = (playerEvals[playerEvals.length-1] ? opponentEvals[opponentEvals.length-1]/playerEvals[playerEvals.length-1] : null);
                        opponentFactors[rowIndex-1][columnIndex-1] = (opponentEvals[opponentEvals.length-1] ? playerEvals[playerEvals.length-1]/opponentEvals[opponentEvals.length-1] : null);

                    } else {
                        gameState[rowIndex-1][columnIndex-1] = gameBoardIn[rowIndex][columnIndex];
                    }
                }
            }

            defensiveMoves = setupMovesReturn(getDefensiveMoves(playerIn, gameState, playerFactors, opponentFactors));

            return defensiveMoves;
        }

    </script>
`

    pagePartHeadClose := "</head>"

    pagePartBodyOpen := `
  <body onload="startGame();">

    <div>
        <table align="center" style="background-color: #004400">
            <tr>
                <td>
`

    pagePartSVGOpen := `
                    <svg
                        id="svgBoard"
                        xmlns="http://www.w3.org/2000/svg"
                        xmlns:xlink="http://www.w3.org/1999/xlink"
                        version="1.1"
                        width="320"
                        height="550"
                        style="background: #004400;">
`
    pagePartSVGBorder := `
                            <polygon
                                points="
                                    %0.10f,%0.10f   
                                    %0.10f,%0.10f   
                                    %0.10f,%0.10f   
                                    %0.10f,%0.10f"
                                style="fill: %s; stroke: black;" />
`
    pagePartSVGLabels := `
                            <text
                                x="%0.10f"
                                y="%0.10f"
                                text-anchor="middle"
                                dominant-baseline="middle"
                                fill="%s"
                                style="font-size:9; font-family: verdana; font-weight: bold;">%s%s</text>
`
    pagePartSVGCellOpen := `
                            <polygon
                                points="
`
    pagePartSVGCellClose := `
                                "
                                onclick="makeCellMove(%d, %d, player);"
                                style="stroke: black; fill: white;"
                            />
`

    pagePartBodyClose := `
                    </svg>
                </td>
                <td style="text-align: center;">
                <table>
                    <tr>
                        <td style="text-align: right;">
                                <span style="color: white;">%s</span>
                        </td>
                        <td colspan="2">
                            <div style="position: relative;">
                                <div id="loader"></div>
                                <img id="opponentFlipCardImg" src="%s/blank.png" width="54" height="72">
                                <img id="opponentPlayCardImg" src="%s/blank.png" width="54" height="72"><br />
                            </div>
                        </td>
                        <td></td>
                    </tr>
                    <tr>
                        <td style="text-align: right;">
                            <span style="color: white;">%s</span>
                        </td>
                        <td colspan="2">
                            <img id="playerFlipCardImg" src="%s/blank.png" width="54" height="72">
                            <img id="playerPlayCardImg" src="%s/blank.png" width="54" height="72"><br />
                        </td>
                        <td>
                            <a style="color: white;" href="http://www.thegamecrafter.com/games/battle-hex">Link to instructions.</a><br />
                            <input id="newGame" type="button" value="New Game" onclick="startGame();" /><br />
                        </td>
                    </tr>
                </table>
                </td>
            </tr>
        </table>
    </div>
  </body>
</html>
`
    if player == helper.FirstPlayer {
        xCoordStr = fmt.Sprintf("%0.10f+((column-row)*%0.10f)", startXCoord, helper.GetCellHalfWidth(cellRadius))
    } else if player == helper.SecondPlayer {
        xCoordStr = fmt.Sprintf("%0.10f+(((%0.10f-column)-(%0.10f-row))*%0.10f)", startXCoord, boardColumns, boardRows, helper.GetCellHalfWidth(cellRadius))
    }

    if player == helper.FirstPlayer {
        yCoordStr = fmt.Sprintf("%0.10f+((column-row)*%0.10f)+((row-1)*%0.10f)", startYCoord, helper.GetCellExtendedLength(cellRadius), 2*helper.GetCellExtendedLength(cellRadius))
    } else if player == helper.SecondPlayer {
        yCoordStr = fmt.Sprintf("%0.10f+(((%0.10f-column)-(%0.10f-row))*%0.10f)+(((%0.10f-row))*%0.10f)", startYCoord, boardColumns+1, boardRows+1, helper.GetCellExtendedLength(cellRadius), boardRows, 2*helper.GetCellExtendedLength(cellRadius))
    }

    if r.URL.Path != "/battlehex_vs_js_ai_v1.1" {
        http.NotFound(w, r)
        return
    }

    fmt.Fprint(w, pagePartStart)
    fmt.Fprint(w, pagePartHeadOpen)
    fmt.Fprintf(w, pagePartJSGameCode, cardset, player, helper.GetPlayerSuits(player), helper.GetOpponentSuits(player), xCoordStr, yCoordStr)
    fmt.Fprint(w, pagePartHeadClose)
    fmt.Fprint(w, pagePartBodyOpen)
    fmt.Fprint(w, pagePartSVGOpen)

    fmt.Fprintf(w, pagePartSVGBorder, startXCoord+helper.GetObtuseBorderXCoord(cellRadius, boardColumns), startYCoord+helper.GetObtuseBorderYCoord(cellRadius, boardColumns), startXCoord, startYCoord+helper.GetAcuteBorderYCoord(cellRadius, boardRows), startXCoord, startYCoord+helper.GetAcuteBorderYCoord(cellRadius, boardRows)+helper.GetAcuteBorderLength(cellRadius), startXCoord+helper.GetObtuseBorderXCoord(cellRadius, boardColumns)+helper.GetObtuseBorderLength(cellRadius), startYCoord+helper.GetObtuseBorderYCoord(cellRadius, boardColumns), "red")
    fmt.Fprintf(w, pagePartSVGBorder, startXCoord-helper.GetObtuseBorderXCoord(cellRadius, boardColumns), startYCoord+helper.GetObtuseBorderYCoord(cellRadius, boardColumns), startXCoord, startYCoord, startXCoord, startYCoord-helper.GetAcuteBorderLength(cellRadius), startXCoord-helper.GetObtuseBorderXCoord(cellRadius, boardColumns)-helper.GetObtuseBorderLength(cellRadius), startYCoord+helper.GetObtuseBorderYCoord(cellRadius, boardColumns), "red")
    fmt.Fprintf(w, pagePartSVGBorder, startXCoord+helper.GetObtuseBorderXCoord(cellRadius, boardColumns), startYCoord+helper.GetObtuseBorderYCoord(cellRadius, boardColumns), startXCoord, startYCoord, startXCoord, startYCoord-helper.GetAcuteBorderLength(cellRadius), startXCoord+helper.GetObtuseBorderXCoord(cellRadius, boardColumns)+helper.GetObtuseBorderLength(cellRadius), startYCoord+helper.GetObtuseBorderYCoord(cellRadius, boardColumns), "black")
    fmt.Fprintf(w, pagePartSVGBorder, startXCoord-helper.GetObtuseBorderXCoord(cellRadius, boardColumns), startYCoord+helper.GetObtuseBorderYCoord(cellRadius, boardColumns), startXCoord, startYCoord+helper.GetAcuteBorderYCoord(cellRadius, boardRows), startXCoord, startYCoord+helper.GetAcuteBorderYCoord(cellRadius, boardRows)+helper.GetAcuteBorderLength(cellRadius), startXCoord-helper.GetObtuseBorderXCoord(cellRadius, boardColumns)-helper.GetObtuseBorderLength(cellRadius), startYCoord+helper.GetObtuseBorderYCoord(cellRadius, boardColumns), "black")

    for i:= 1; i <= boardRows; i++ {
        fmt.Fprintf(w, pagePartSVGLabels, helper.GetCellXCoord(cellRadius, startXCoord, boardRows, boardColumns, i, 0, boardShape+player), helper.GetCellYCoord(cellRadius, startYCoord, boardRows, boardColumns, i, 0, boardShape+player), "black", "R", helper.GetBoardLabel(i))
        fmt.Fprintf(w, pagePartSVGLabels, helper.GetCellXCoord(cellRadius, startXCoord, boardRows, boardColumns, i, 14, boardShape+player), helper.GetCellYCoord(cellRadius, startYCoord, boardRows, boardColumns, i, 14, boardShape+player), "black", "R", helper.GetBoardLabel(i))
    }
    for i:= 1; i <= boardColumns; i++ {
        fmt.Fprintf(w, pagePartSVGLabels, helper.GetCellXCoord(cellRadius, startXCoord, boardRows, boardColumns, 0, i, boardShape+player), helper.GetCellYCoord(cellRadius, startYCoord, boardRows, boardColumns, 0, i, boardShape+player), "white", "B", helper.GetBoardLabel(i))
        fmt.Fprintf(w, pagePartSVGLabels, helper.GetCellXCoord(cellRadius, startXCoord, boardRows, boardColumns, 14, i, boardShape+player), helper.GetCellYCoord(cellRadius, startYCoord, boardRows, boardColumns, 14, i, boardShape+player), "white", "B", helper.GetBoardLabel(i))
    }

    for i:= 1; i <= boardRows; i++ {
        for j:= 1; j <= boardColumns; j++ {
            fmt.Fprint(w, pagePartSVGCellOpen)
            for k:= 0; k < 6; k++ {
                fmt.Fprintf(w, "                                    %0.10f,%0.10f\n", helper.GetPointXCoord(cellRadius, helper.GetCellXCoord(cellRadius, startXCoord, boardRows, boardColumns, i, j, boardShape+player), k, boardShape+player), helper.GetPointYCoord(cellRadius, helper.GetCellYCoord(cellRadius, startYCoord, boardRows, boardColumns, i, j, boardShape+player), k, boardShape+player))
            }
            fmt.Fprintf(w, pagePartSVGCellClose, i, j)
        }
    }
            
    fmt.Fprintf(w, pagePartBodyClose, opponentMatch, cardset, cardset, playerMatch, cardset, cardset)
}
