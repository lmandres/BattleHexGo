package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "battle-hex-go/helper"
)

func main() {
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/battlehex_vs_js_ai_v1.1", battleHexJSHandler)

    port :=  os.Getenv("PORT")
    if port == "" {
        port = "8080"
        log.Printf("Defaulting to port %s", port)
    }

    log.Printf("Listening on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    fmt.Fprint(w, "Hello, World!")
}

func battleHexJSHandler(w http.ResponseWriter, r *http.Request) {

    var cardset string = "/static/cardstux/"

    var player int = helper.FirstPlayer
    var imageWidth int = 320
    var imageHeight int = 550
    var boardRows int = 13
    var boardColumns int = 13
    var boardShape int = helper.VerticalBoard

    var xCoordStr string = ""
    var yCoordStr string = ""

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
					document.getElementById("opponentPlayCardImg").src = cardset + "back.png";
					document.getElementById("opponentFlipCardImg").src = cardset + "back.png";
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

					document.getElementById("playerFlipCardImg").src = cardset + playerSuits["b,r".split(",").indexOf(playerFlipCard.substr(0, 1))] + playerFlipCard.substr(1) + ".png";
					document.getElementById("playerPlayCardImg").src = cardset + playerSuits["b,r".split(",").indexOf(playerPlayCard.substr(0, 1))] + playerPlayCard.substr(1) + ".png";
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

				document.getElementById("opponentFlipCardImg").src = cardset + computerSuits["b,r".split(",").indexOf(computerFlipCard.substr(0, 1))] + computerFlipCard.substr(1) + ".png";
				document.getElementById("opponentPlayCardImg").src = cardset + computerSuits["b,r".split(",").indexOf(computerPlayCard.substr(0, 1))] + computerPlayCard.substr(1) + ".png";

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

									document.getElementById("opponentFlipCardImg").src = cardset + computerSuits["b,r".split(",").indexOf(computerFlipCard.substr(0, 1))] + computerFlipCard.substr(1) + ".png";
									document.getElementById("opponentPlayCardImg").src = cardset + computerSuits["b,r".split(",").indexOf(computerPlayCard.substr(0, 1))] + computerPlayCard.substr(1) + ".png";

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

									document.getElementById("opponentFlipCardImg").src = cardset + computerSuits["b,r".split(",").indexOf(computerFlipCard.substr(0, 1))] + computerFlipCard.substr(1) + ".png";
									document.getElementById("opponentPlayCardImg").src = cardset + computerSuits["b,r".split(",").indexOf(computerPlayCard.substr(0, 1))] + computerPlayCard.substr(1) + ".png";

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
				
				document.getElementById("playerFlipCardImg").src = cardset + playerSuits["b,r".split(",").indexOf(playerFlipCard.substr(0, 1))] + playerFlipCard.substr(1) + ".png";
				document.getElementById("playerPlayCardImg").src = cardset + playerSuits["b,r".split(",").indexOf(playerPlayCard.substr(0, 1))] + playerPlayCard.substr(1) + ".png";

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

			document.getElementById("playerPlayCardImg").src = cardset + "blank.png";
			document.getElementById("playerFlipCardImg").src = cardset + "blank.png";
			document.getElementById("opponentPlayCardImg").src = cardset + "blank.png";
			document.getElementById("opponentFlipCardImg").src = cardset + "blank.png";

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

    pagePartSVGLabels2 := `
							<text
								x="139.51337754488424"
								y="62.096774193548384"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R2</text>
							<text
								x="282.91973473069453"
								y="310.48387096774195"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R2</text>

							<text
								x="129.27006631732638"
								y="79.83870967741936"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R3</text>
							<text
								x="272.6764235031367"
								y="328.2258064516129"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R3</text>

							<text
								x="119.0267550897685"
								y="97.58064516129032"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R4</text>
							<text
								x="262.43311227557876"
								y="345.9677419354839"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R4</text>

							<text
								x="108.78344386221062"
								y="115.32258064516128"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R5</text>
							<text
								x="252.1898010480209"
								y="363.7096774193549"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R5</text>

							<text
								x="98.54013263465274"
								y="133.06451612903226"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R6</text>
							<text
								x="241.946489820463"
								y="381.45161290322585"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R6</text>

							<text
								x="88.29682140709487"
								y="150.80645161290323"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R7</text>
							<text
								x="231.70317859290515"
								y="399.19354838709677"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R7</text>

							<text
								x="78.05351017953699"
								y="168.54838709677418"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R8</text>
							<text
								x="221.45986736534726"
								y="416.93548387096774"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R8</text>

							<text
								x="67.8101989519791"
								y="186.29032258064515"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R9</text>
							<text
								x="211.21655613778938"
								y="434.6774193548387"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R9</text>

							<text
								x="57.566887724421235"
								y="204.03225806451613"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R10</text>
							<text
								x="200.9732449102315"
								y="452.4193548387097"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">R10</text>

							<text
								x="47.32357649686335"
								y="221.77419354838713"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">RJ</text>
							<text
								x="190.72993368267362"
								y="470.16129032258067"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">RJ</text>

							<text
								x="37.08026526930547"
								y="239.51612903225805"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">RQ</text>
							<text
								x="180.48662245511576"
								y="487.9032258064516"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">RQ</text>

							<text
								x="26.83695404174759"
								y="257.258064516129"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">RK</text>
							<text
								x="170.24331122755788"
								y="505.64516129032256"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="black"
								style="font-size:9; font-family: verdana; font-weight: bold;">RK</text>



							<text
								x="170.24331122755788"
								y="44.35483870967742"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">BA</text>
							<text
								x="26.83695404174759"
								y="292.741935483871"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">BA</text>

							<text
								x="180.48662245511576"
								y="62.096774193548384"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B2</text>
							<text
								x="37.08026526930547"
								y="310.48387096774195"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B2</text>

							<text
								x="190.72993368267362"
								y="79.83870967741935"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B3</text>
							<text
								x="47.32357649686335"
								y="328.2258064516129"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B3</text>

							<text
								x="200.9732449102315"
								y="97.58064516129032"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B4</text>
							<text
								x="57.566887724421235"
								y="345.9677419354839"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B4</text>

							<text
								x="211.21655613778938"
								y="115.3225806451613"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B5</text>
							<text
								x="67.8101989519791"
								y="363.7096774193549"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B5</text>

							<text
								x="221.45986736534726"
								y="133.06451612903226"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B6</text>
							<text
								x="78.05351017953699"
								y="381.45161290322585"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B6</text>

							<text
								x="231.70317859290515"
								y="150.80645161290323"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B7</text>
							<text
								x="88.29682140709487"
								y="399.1935483870968"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B7</text>

							<text
								x="241.946489820463"
								y="168.5483870967742"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B8</text>
							<text
								x="98.54013263465274"
								y="416.93548387096774"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B8</text>

							<text
								x="252.1898010480209"
								y="186.29032258064518"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B9</text>
							<text
								x="108.78344386221062"
								y="434.6774193548387"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B9</text>

							<text
								x="262.43311227557876"
								y="204.03225806451616"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B10</text>
							<text
								x="119.0267550897685"
								y="452.4193548387097"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">B10</text>

							<text
								x="272.6764235031367"
								y="221.77419354838707"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">BJ</text>
							<text
								x="129.27006631732638"
								y="470.16129032258067"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">BJ</text>

							<text
								x="282.91973473069453"
								y="239.51612903225805"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">BQ</text>
							<text
								x="139.51337754488424"
								y="487.90322580645164"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">BQ</text>

							<text
								x="293.1630459582524"
								y="257.258064516129"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">BK</text>
							<text
								x="149.75668877244212"
								y="505.6451612903226"
								text-anchor="middle"
								dominant-baseline="middle"
								fill="white"
								style="font-size:9; font-family: verdana; font-weight: bold;">BK</text>



							<polygon
								points="

									170.24331122755788,68.01075268817203

									160.0,73.9247311827957

									149.75668877244212,68.01075268817205

									149.75668877244212,56.182795698924735

									160.0,50.26881720430107

									170.24331122755788,56.182795698924735

								"
								onclick="makeCellMove(1, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,85.75268817204301

									170.24331122755788,91.66666666666667

									160.0,85.75268817204302

									160.0,73.92473118279571

									170.24331122755788,68.01075268817205

									180.48662245511576,73.92473118279571

								"
								onclick="makeCellMove(1, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,103.49462365591397

									180.48662245511576,109.40860215053763

									170.24331122755788,103.49462365591398

									170.24331122755788,91.66666666666667

									180.48662245511576,85.75268817204301

									190.72993368267365,91.66666666666667

								"
								onclick="makeCellMove(1, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,121.23655913978493

									190.72993368267362,127.1505376344086

									180.48662245511574,121.23655913978494

									180.48662245511574,109.40860215053763

									190.72993368267362,103.49462365591397

									200.9732449102315,109.40860215053763

								"
								onclick="makeCellMove(1, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,138.9784946236559

									200.9732449102315,144.89247311827955

									190.72993368267362,138.9784946236559

									190.72993368267362,127.15053763440861

									200.9732449102315,121.23655913978494

									211.21655613778938,127.15053763440861

								"
								onclick="makeCellMove(1, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									221.45986736534726,156.72043010752688

									211.21655613778938,162.63440860215053

									200.9732449102315,156.72043010752688

									200.9732449102315,144.89247311827958

									211.21655613778938,138.97849462365593

									221.45986736534726,144.89247311827958

								"
								onclick="makeCellMove(1, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									231.70317859290515,174.46236559139786

									221.45986736534726,180.3763440860215

									211.21655613778938,174.46236559139786

									211.21655613778938,162.63440860215056

									221.45986736534726,156.7204301075269

									231.70317859290515,162.63440860215056

								"
								onclick="makeCellMove(1, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									241.94648982046303,192.2043010752688

									231.70317859290515,198.11827956989248

									221.45986736534726,192.2043010752688

									221.45986736534726,180.3763440860215

									231.70317859290515,174.46236559139783

									241.94648982046303,180.3763440860215

								"
								onclick="makeCellMove(1, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									252.18980104802088,209.94623655913978

									241.946489820463,215.86021505376345

									231.70317859290512,209.94623655913978

									231.70317859290512,198.11827956989248

									241.946489820463,192.2043010752688

									252.18980104802088,198.11827956989248

								"
								onclick="makeCellMove(1, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									262.43311227557876,227.68817204301075

									252.1898010480209,233.60215053763443

									241.94648982046303,227.68817204301075

									241.94648982046303,215.86021505376345

									252.1898010480209,209.94623655913978

									262.43311227557876,215.86021505376345

								"
								onclick="makeCellMove(1, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									272.6764235031366,245.43010752688173

									262.43311227557876,251.3440860215054

									252.18980104802088,245.43010752688173

									252.18980104802088,233.60215053763443

									262.43311227557876,227.68817204301075

									272.6764235031366,233.60215053763443

								"
								onclick="makeCellMove(1, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									282.91973473069453,263.1720430107527

									272.6764235031367,269.0860215053763

									262.4331122755788,263.1720430107527

									262.4331122755788,251.34408602150538

									272.6764235031367,245.43010752688173

									282.91973473069453,251.34408602150538

								"
								onclick="makeCellMove(1, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									293.1630459582524,280.9139784946237

									282.91973473069453,286.8279569892473

									272.6764235031367,280.9139784946237

									272.6764235031367,269.0860215053763

									282.91973473069453,263.1720430107527

									293.1630459582524,269.0860215053763

								"
								onclick="makeCellMove(1, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,85.75268817204301

									149.75668877244212,91.66666666666667

									139.51337754488424,85.75268817204302

									139.51337754488424,73.92473118279571

									149.75668877244212,68.01075268817205

									160.0,73.92473118279571

								"
								onclick="makeCellMove(2, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,103.49462365591397

									160.0,109.40860215053763

									149.75668877244212,103.49462365591398

									149.75668877244212,91.66666666666667

									160.0,85.75268817204301

									170.24331122755788,91.66666666666667

								"
								onclick="makeCellMove(2, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,121.23655913978494

									170.24331122755788,127.15053763440861

									160.0,121.23655913978496

									160.0,109.40860215053765

									170.24331122755788,103.49462365591398

									180.48662245511576,109.40860215053765

								"
								onclick="makeCellMove(2, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,138.9784946236559

									180.48662245511576,144.89247311827955

									170.24331122755788,138.9784946236559

									170.24331122755788,127.15053763440861

									180.48662245511576,121.23655913978494

									190.72993368267365,127.15053763440861

								"
								onclick="makeCellMove(2, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,156.72043010752688

									190.72993368267362,162.63440860215053

									180.48662245511574,156.72043010752688

									180.48662245511574,144.89247311827958

									190.72993368267362,138.97849462365593

									200.9732449102315,144.89247311827958

								"
								onclick="makeCellMove(2, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,174.46236559139786

									200.9732449102315,180.3763440860215

									190.72993368267362,174.46236559139786

									190.72993368267362,162.63440860215056

									200.9732449102315,156.7204301075269

									211.21655613778938,162.63440860215056

								"
								onclick="makeCellMove(2, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									221.45986736534726,192.20430107526883

									211.21655613778938,198.11827956989248

									200.9732449102315,192.20430107526883

									200.9732449102315,180.37634408602153

									211.21655613778938,174.46236559139788

									221.45986736534726,180.37634408602153

								"
								onclick="makeCellMove(2, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									231.70317859290515,209.9462365591398

									221.45986736534726,215.86021505376345

									211.21655613778938,209.9462365591398

									211.21655613778938,198.1182795698925

									221.45986736534726,192.20430107526886

									231.70317859290515,198.1182795698925

								"
								onclick="makeCellMove(2, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									241.94648982046303,227.68817204301072

									231.70317859290515,233.60215053763437

									221.45986736534726,227.68817204301072

									221.45986736534726,215.86021505376343

									231.70317859290515,209.94623655913978

									241.94648982046303,215.86021505376343

								"
								onclick="makeCellMove(2, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									252.18980104802088,245.4301075268817

									241.946489820463,251.34408602150535

									231.70317859290512,245.4301075268817

									231.70317859290512,233.6021505376344

									241.946489820463,227.68817204301075

									252.18980104802088,233.6021505376344

								"
								onclick="makeCellMove(2, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									262.43311227557876,263.1720430107527

									252.1898010480209,269.0860215053763

									241.94648982046303,263.1720430107527

									241.94648982046303,251.34408602150538

									252.1898010480209,245.43010752688173

									262.43311227557876,251.34408602150538

								"
								onclick="makeCellMove(2, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									272.6764235031366,280.9139784946237

									262.43311227557876,286.8279569892473

									252.18980104802088,280.9139784946237

									252.18980104802088,269.0860215053763

									262.43311227557876,263.1720430107527

									272.6764235031366,269.0860215053763

								"
								onclick="makeCellMove(2, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									282.91973473069453,298.65591397849465

									272.6764235031367,304.5698924731183

									262.4331122755788,298.65591397849465

									262.4331122755788,286.8279569892473

									272.6764235031367,280.9139784946237

									282.91973473069453,286.8279569892473

								"
								onclick="makeCellMove(2, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,103.49462365591397

									139.51337754488424,109.40860215053763

									129.27006631732635,103.49462365591398

									129.27006631732635,91.66666666666667

									139.51337754488424,85.75268817204301

									149.75668877244212,91.66666666666667

								"
								onclick="makeCellMove(3, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,121.23655913978493

									149.75668877244212,127.1505376344086

									139.51337754488424,121.23655913978494

									139.51337754488424,109.40860215053763

									149.75668877244212,103.49462365591397

									160.0,109.40860215053763

								"
								onclick="makeCellMove(3, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,138.9784946236559

									160.0,144.89247311827955

									149.75668877244212,138.9784946236559

									149.75668877244212,127.15053763440861

									160.0,121.23655913978494

									170.24331122755788,127.15053763440861

								"
								onclick="makeCellMove(3, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,156.72043010752688

									170.24331122755788,162.63440860215053

									160.0,156.72043010752688

									160.0,144.89247311827958

									170.24331122755788,138.97849462365593

									180.48662245511576,144.89247311827958

								"
								onclick="makeCellMove(3, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,174.46236559139786

									180.48662245511576,180.3763440860215

									170.24331122755788,174.46236559139786

									170.24331122755788,162.63440860215056

									180.48662245511576,156.7204301075269

									190.72993368267365,162.63440860215056

								"
								onclick="makeCellMove(3, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,192.2043010752688

									190.72993368267362,198.11827956989248

									180.48662245511574,192.2043010752688

									180.48662245511574,180.3763440860215

									190.72993368267362,174.46236559139783

									200.9732449102315,180.3763440860215

								"
								onclick="makeCellMove(3, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,209.94623655913978

									200.9732449102315,215.86021505376345

									190.72993368267362,209.94623655913978

									190.72993368267362,198.11827956989248

									200.9732449102315,192.2043010752688

									211.21655613778938,198.11827956989248

								"
								onclick="makeCellMove(3, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									221.45986736534726,227.68817204301075

									211.21655613778938,233.60215053763443

									200.9732449102315,227.68817204301075

									200.9732449102315,215.86021505376345

									211.21655613778938,209.94623655913978

									221.45986736534726,215.86021505376345

								"
								onclick="makeCellMove(3, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									231.70317859290515,245.43010752688173

									221.45986736534726,251.3440860215054

									211.21655613778938,245.43010752688173

									211.21655613778938,233.60215053763443

									221.45986736534726,227.68817204301075

									231.70317859290515,233.60215053763443

								"
								onclick="makeCellMove(3, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									241.94648982046303,263.1720430107527

									231.70317859290515,269.0860215053763

									221.45986736534726,263.1720430107527

									221.45986736534726,251.34408602150538

									231.70317859290515,245.43010752688173

									241.94648982046303,251.34408602150538

								"
								onclick="makeCellMove(3, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									252.18980104802088,280.9139784946237

									241.946489820463,286.8279569892473

									231.70317859290512,280.9139784946237

									231.70317859290512,269.0860215053763

									241.946489820463,263.1720430107527

									252.18980104802088,269.0860215053763

								"
								onclick="makeCellMove(3, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									262.43311227557876,298.65591397849465

									252.1898010480209,304.5698924731183

									241.94648982046303,298.65591397849465

									241.94648982046303,286.8279569892473

									252.1898010480209,280.9139784946237

									262.43311227557876,286.8279569892473

								"
								onclick="makeCellMove(3, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									272.6764235031366,316.3978494623656

									262.43311227557876,322.31182795698925

									252.18980104802088,316.3978494623656

									252.18980104802088,304.5698924731183

									262.43311227557876,298.65591397849465

									272.6764235031366,304.5698924731183

								"
								onclick="makeCellMove(3, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,121.23655913978493

									129.27006631732638,127.1505376344086

									119.0267550897685,121.23655913978494

									119.0267550897685,109.40860215053763

									129.27006631732638,103.49462365591397

									139.51337754488426,109.40860215053763

								"
								onclick="makeCellMove(4, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,138.9784946236559

									139.51337754488424,144.89247311827955

									129.27006631732635,138.9784946236559

									129.27006631732635,127.15053763440861

									139.51337754488424,121.23655913978494

									149.75668877244212,127.15053763440861

								"
								onclick="makeCellMove(4, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,156.72043010752688

									149.75668877244212,162.63440860215053

									139.51337754488424,156.72043010752688

									139.51337754488424,144.89247311827958

									149.75668877244212,138.97849462365593

									160.0,144.89247311827958

								"
								onclick="makeCellMove(4, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,174.46236559139786

									160.0,180.3763440860215

									149.75668877244212,174.46236559139786

									149.75668877244212,162.63440860215056

									160.0,156.7204301075269

									170.24331122755788,162.63440860215056

								"
								onclick="makeCellMove(4, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,192.20430107526883

									170.24331122755788,198.11827956989248

									160.0,192.20430107526883

									160.0,180.37634408602153

									170.24331122755788,174.46236559139788

									180.48662245511576,180.37634408602153

								"
								onclick="makeCellMove(4, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,209.94623655913978

									180.48662245511576,215.86021505376345

									170.24331122755788,209.94623655913978

									170.24331122755788,198.11827956989248

									180.48662245511576,192.2043010752688

									190.72993368267365,198.11827956989248

								"
								onclick="makeCellMove(4, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,227.68817204301072

									190.72993368267362,233.60215053763437

									180.48662245511574,227.68817204301072

									180.48662245511574,215.86021505376343

									190.72993368267362,209.94623655913978

									200.9732449102315,215.86021505376343

								"
								onclick="makeCellMove(4, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,245.4301075268817

									200.9732449102315,251.34408602150535

									190.72993368267362,245.4301075268817

									190.72993368267362,233.6021505376344

									200.9732449102315,227.68817204301075

									211.21655613778938,233.6021505376344

								"
								onclick="makeCellMove(4, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									221.45986736534726,263.1720430107527

									211.21655613778938,269.0860215053763

									200.9732449102315,263.1720430107527

									200.9732449102315,251.34408602150538

									211.21655613778938,245.43010752688173

									221.45986736534726,251.34408602150538

								"
								onclick="makeCellMove(4, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									231.70317859290515,280.9139784946237

									221.45986736534726,286.8279569892473

									211.21655613778938,280.9139784946237

									211.21655613778938,269.0860215053763

									221.45986736534726,263.1720430107527

									231.70317859290515,269.0860215053763

								"
								onclick="makeCellMove(4, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									241.94648982046303,298.65591397849465

									231.70317859290515,304.5698924731183

									221.45986736534726,298.65591397849465

									221.45986736534726,286.8279569892473

									231.70317859290515,280.9139784946237

									241.94648982046303,286.8279569892473

								"
								onclick="makeCellMove(4, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									252.18980104802088,316.3978494623656

									241.946489820463,322.31182795698925

									231.70317859290512,316.3978494623656

									231.70317859290512,304.5698924731183

									241.946489820463,298.65591397849465

									252.18980104802088,304.5698924731183

								"
								onclick="makeCellMove(4, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									262.43311227557876,334.1397849462366

									252.1898010480209,340.0537634408602

									241.94648982046303,334.1397849462366

									241.94648982046303,322.31182795698925

									252.1898010480209,316.3978494623656

									262.43311227557876,322.31182795698925

								"
								onclick="makeCellMove(4, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,138.9784946236559

									119.0267550897685,144.89247311827955

									108.78344386221062,138.9784946236559

									108.78344386221062,127.15053763440861

									119.0267550897685,121.23655913978494

									129.27006631732638,127.15053763440861

								"
								onclick="makeCellMove(5, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,156.72043010752688

									129.27006631732638,162.63440860215053

									119.0267550897685,156.72043010752688

									119.0267550897685,144.89247311827958

									129.27006631732638,138.97849462365593

									139.51337754488426,144.89247311827958

								"
								onclick="makeCellMove(5, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,174.46236559139786

									139.51337754488424,180.3763440860215

									129.27006631732635,174.46236559139786

									129.27006631732635,162.63440860215056

									139.51337754488424,156.7204301075269

									149.75668877244212,162.63440860215056

								"
								onclick="makeCellMove(5, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,192.2043010752688

									149.75668877244212,198.11827956989248

									139.51337754488424,192.2043010752688

									139.51337754488424,180.3763440860215

									149.75668877244212,174.46236559139783

									160.0,180.3763440860215

								"
								onclick="makeCellMove(5, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,209.94623655913978

									160.0,215.86021505376345

									149.75668877244212,209.94623655913978

									149.75668877244212,198.11827956989248

									160.0,192.2043010752688

									170.24331122755788,198.11827956989248

								"
								onclick="makeCellMove(5, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,227.68817204301075

									170.24331122755788,233.60215053763443

									160.0,227.68817204301075

									160.0,215.86021505376345

									170.24331122755788,209.94623655913978

									180.48662245511576,215.86021505376345

								"
								onclick="makeCellMove(5, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,245.4301075268817

									180.48662245511576,251.34408602150535

									170.24331122755788,245.4301075268817

									170.24331122755788,233.6021505376344

									180.48662245511576,227.68817204301075

									190.72993368267365,233.6021505376344

								"
								onclick="makeCellMove(5, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,263.1720430107527

									190.72993368267362,269.0860215053763

									180.48662245511574,263.1720430107527

									180.48662245511574,251.34408602150538

									190.72993368267362,245.43010752688173

									200.9732449102315,251.34408602150538

								"
								onclick="makeCellMove(5, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,280.9139784946237

									200.9732449102315,286.8279569892473

									190.72993368267362,280.9139784946237

									190.72993368267362,269.0860215053763

									200.9732449102315,263.1720430107527

									211.21655613778938,269.0860215053763

								"
								onclick="makeCellMove(5, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									221.45986736534726,298.65591397849465

									211.21655613778938,304.5698924731183

									200.9732449102315,298.65591397849465

									200.9732449102315,286.8279569892473

									211.21655613778938,280.9139784946237

									221.45986736534726,286.8279569892473

								"
								onclick="makeCellMove(5, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									231.70317859290515,316.3978494623656

									221.45986736534726,322.31182795698925

									211.21655613778938,316.3978494623656

									211.21655613778938,304.5698924731183

									221.45986736534726,298.65591397849465

									231.70317859290515,304.5698924731183

								"
								onclick="makeCellMove(5, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									241.94648982046303,334.1397849462366

									231.70317859290515,340.0537634408602

									221.45986736534726,334.1397849462366

									221.45986736534726,322.31182795698925

									231.70317859290515,316.3978494623656

									241.94648982046303,322.31182795698925

								"
								onclick="makeCellMove(5, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									252.18980104802088,351.8817204301076

									241.946489820463,357.7956989247312

									231.70317859290512,351.8817204301076

									231.70317859290512,340.0537634408602

									241.946489820463,334.1397849462366

									252.18980104802088,340.0537634408602

								"
								onclick="makeCellMove(5, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									119.0267550897685,156.72043010752688

									108.78344386221062,162.63440860215053

									98.54013263465274,156.72043010752688

									98.54013263465274,144.89247311827958

									108.78344386221062,138.97849462365593

									119.0267550897685,144.89247311827958

								"
								onclick="makeCellMove(6, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,174.46236559139786

									119.0267550897685,180.3763440860215

									108.78344386221062,174.46236559139786

									108.78344386221062,162.63440860215056

									119.0267550897685,156.7204301075269

									129.27006631732638,162.63440860215056

								"
								onclick="makeCellMove(6, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,192.20430107526883

									129.27006631732638,198.11827956989248

									119.0267550897685,192.20430107526883

									119.0267550897685,180.37634408602153

									129.27006631732638,174.46236559139788

									139.51337754488426,180.37634408602153

								"
								onclick="makeCellMove(6, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,209.9462365591398

									139.51337754488424,215.86021505376345

									129.27006631732635,209.9462365591398

									129.27006631732635,198.1182795698925

									139.51337754488424,192.20430107526886

									149.75668877244212,198.1182795698925

								"
								onclick="makeCellMove(6, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,227.68817204301075

									149.75668877244212,233.60215053763443

									139.51337754488424,227.68817204301075

									139.51337754488424,215.86021505376345

									149.75668877244212,209.94623655913978

									160.0,215.86021505376345

								"
								onclick="makeCellMove(6, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,245.43010752688173

									160.0,251.3440860215054

									149.75668877244212,245.43010752688173

									149.75668877244212,233.60215053763443

									160.0,227.68817204301075

									170.24331122755788,233.60215053763443

								"
								onclick="makeCellMove(6, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,263.1720430107527

									170.24331122755788,269.0860215053763

									160.0,263.1720430107527

									160.0,251.34408602150538

									170.24331122755788,245.43010752688173

									180.48662245511576,251.34408602150538

								"
								onclick="makeCellMove(6, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,280.9139784946237

									180.48662245511576,286.8279569892473

									170.24331122755788,280.9139784946237

									170.24331122755788,269.0860215053763

									180.48662245511576,263.1720430107527

									190.72993368267365,269.0860215053763

								"
								onclick="makeCellMove(6, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,298.65591397849465

									190.72993368267362,304.5698924731183

									180.48662245511574,298.65591397849465

									180.48662245511574,286.8279569892473

									190.72993368267362,280.9139784946237

									200.9732449102315,286.8279569892473

								"
								onclick="makeCellMove(6, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,316.3978494623656

									200.9732449102315,322.31182795698925

									190.72993368267362,316.3978494623656

									190.72993368267362,304.5698924731183

									200.9732449102315,298.65591397849465

									211.21655613778938,304.5698924731183

								"
								onclick="makeCellMove(6, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									221.45986736534726,334.1397849462366

									211.21655613778938,340.0537634408602

									200.9732449102315,334.1397849462366

									200.9732449102315,322.31182795698925

									211.21655613778938,316.3978494623656

									221.45986736534726,322.31182795698925

								"
								onclick="makeCellMove(6, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									231.70317859290515,351.8817204301076

									221.45986736534726,357.7956989247312

									211.21655613778938,351.8817204301076

									211.21655613778938,340.0537634408602

									221.45986736534726,334.1397849462366

									231.70317859290515,340.0537634408602

								"
								onclick="makeCellMove(6, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									241.94648982046303,369.62365591397855

									231.70317859290515,375.5376344086022

									221.45986736534726,369.62365591397855

									221.45986736534726,357.7956989247312

									231.70317859290515,351.8817204301076

									241.94648982046303,357.7956989247312

								"
								onclick="makeCellMove(6, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									108.78344386221062,174.46236559139786

									98.54013263465274,180.3763440860215

									88.29682140709485,174.46236559139786

									88.29682140709485,162.63440860215056

									98.54013263465274,156.7204301075269

									108.78344386221062,162.63440860215056

								"
								onclick="makeCellMove(7, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									119.0267550897685,192.2043010752688

									108.78344386221062,198.11827956989248

									98.54013263465274,192.2043010752688

									98.54013263465274,180.3763440860215

									108.78344386221062,174.46236559139783

									119.0267550897685,180.3763440860215

								"
								onclick="makeCellMove(7, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,209.94623655913978

									119.0267550897685,215.86021505376345

									108.78344386221062,209.94623655913978

									108.78344386221062,198.11827956989248

									119.0267550897685,192.2043010752688

									129.27006631732638,198.11827956989248

								"
								onclick="makeCellMove(7, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,227.68817204301075

									129.27006631732638,233.60215053763443

									119.0267550897685,227.68817204301075

									119.0267550897685,215.86021505376345

									129.27006631732638,209.94623655913978

									139.51337754488426,215.86021505376345

								"
								onclick="makeCellMove(7, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,245.4301075268817

									139.51337754488424,251.34408602150535

									129.27006631732635,245.4301075268817

									129.27006631732635,233.6021505376344

									139.51337754488424,227.68817204301075

									149.75668877244212,233.6021505376344

								"
								onclick="makeCellMove(7, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,263.1720430107527

									149.75668877244212,269.0860215053763

									139.51337754488424,263.1720430107527

									139.51337754488424,251.34408602150538

									149.75668877244212,245.43010752688173

									160.0,251.34408602150538

								"
								onclick="makeCellMove(7, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,280.9139784946237

									160.0,286.8279569892473

									149.75668877244212,280.9139784946237

									149.75668877244212,269.0860215053763

									160.0,263.1720430107527

									170.24331122755788,269.0860215053763

								"
								onclick="makeCellMove(7, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,298.65591397849465

									170.24331122755788,304.5698924731183

									160.0,298.65591397849465

									160.0,286.8279569892473

									170.24331122755788,280.9139784946237

									180.48662245511576,286.8279569892473

								"
								onclick="makeCellMove(7, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,316.3978494623656

									180.48662245511576,322.31182795698925

									170.24331122755788,316.3978494623656

									170.24331122755788,304.5698924731183

									180.48662245511576,298.65591397849465

									190.72993368267365,304.5698924731183

								"
								onclick="makeCellMove(7, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,334.1397849462366

									190.72993368267362,340.0537634408602

									180.48662245511574,334.1397849462366

									180.48662245511574,322.31182795698925

									190.72993368267362,316.3978494623656

									200.9732449102315,322.31182795698925

								"
								onclick="makeCellMove(7, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,351.8817204301076

									200.9732449102315,357.7956989247312

									190.72993368267362,351.8817204301076

									190.72993368267362,340.0537634408602

									200.9732449102315,334.1397849462366

									211.21655613778938,340.0537634408602

								"
								onclick="makeCellMove(7, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									221.45986736534726,369.62365591397855

									211.21655613778938,375.5376344086022

									200.9732449102315,369.62365591397855

									200.9732449102315,357.7956989247312

									211.21655613778938,351.8817204301076

									221.45986736534726,357.7956989247312

								"
								onclick="makeCellMove(7, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									231.70317859290515,387.3655913978495

									221.45986736534726,393.27956989247315

									211.21655613778938,387.3655913978495

									211.21655613778938,375.5376344086022

									221.45986736534726,369.62365591397855

									231.70317859290515,375.5376344086022

								"
								onclick="makeCellMove(7, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									98.54013263465275,192.2043010752688

									88.29682140709487,198.11827956989248

									78.053510179537,192.2043010752688

									78.05351017953699,180.3763440860215

									88.29682140709487,174.46236559139783

									98.54013263465275,180.3763440860215

								"
								onclick="makeCellMove(8, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									108.78344386221062,209.94623655913975

									98.54013263465274,215.8602150537634

									88.29682140709485,209.94623655913975

									88.29682140709485,198.11827956989245

									98.54013263465274,192.2043010752688

									108.78344386221062,198.11827956989245

								"
								onclick="makeCellMove(8, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									119.0267550897685,227.68817204301072

									108.78344386221062,233.60215053763437

									98.54013263465274,227.68817204301072

									98.54013263465274,215.86021505376343

									108.78344386221062,209.94623655913978

									119.0267550897685,215.86021505376343

								"
								onclick="makeCellMove(8, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,245.4301075268817

									119.0267550897685,251.34408602150535

									108.78344386221062,245.4301075268817

									108.78344386221062,233.6021505376344

									119.0267550897685,227.68817204301075

									129.27006631732638,233.6021505376344

								"
								onclick="makeCellMove(8, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,263.1720430107527

									129.27006631732638,269.0860215053763

									119.0267550897685,263.1720430107527

									119.0267550897685,251.34408602150538

									129.27006631732638,245.43010752688173

									139.51337754488426,251.34408602150538

								"
								onclick="makeCellMove(8, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,280.9139784946237

									139.51337754488424,286.8279569892473

									129.27006631732635,280.9139784946237

									129.27006631732635,269.0860215053763

									139.51337754488424,263.1720430107527

									149.75668877244212,269.0860215053763

								"
								onclick="makeCellMove(8, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,298.65591397849465

									149.75668877244212,304.5698924731183

									139.51337754488424,298.65591397849465

									139.51337754488424,286.8279569892473

									149.75668877244212,280.9139784946237

									160.0,286.8279569892473

								"
								onclick="makeCellMove(8, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,316.3978494623656

									160.0,322.31182795698925

									149.75668877244212,316.3978494623656

									149.75668877244212,304.5698924731183

									160.0,298.65591397849465

									170.24331122755788,304.5698924731183

								"
								onclick="makeCellMove(8, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,334.1397849462366

									170.24331122755788,340.0537634408602

									160.0,334.1397849462366

									160.0,322.31182795698925

									170.24331122755788,316.3978494623656

									180.48662245511576,322.31182795698925

								"
								onclick="makeCellMove(8, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,351.8817204301075

									180.48662245511576,357.79569892473114

									170.24331122755788,351.8817204301075

									170.24331122755788,340.05376344086017

									180.48662245511576,334.13978494623655

									190.72993368267365,340.05376344086017

								"
								onclick="makeCellMove(8, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,369.6236559139785

									190.72993368267362,375.5376344086021

									180.48662245511574,369.6236559139785

									180.48662245511574,357.79569892473114

									190.72993368267362,351.8817204301075

									200.9732449102315,357.79569892473114

								"
								onclick="makeCellMove(8, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,387.36559139784947

									200.9732449102315,393.2795698924731

									190.72993368267362,387.36559139784947

									190.72993368267362,375.5376344086021

									200.9732449102315,369.6236559139785

									211.21655613778938,375.5376344086021

								"
								onclick="makeCellMove(8, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									221.45986736534726,405.10752688172045

									211.21655613778938,411.02150537634407

									200.9732449102315,405.10752688172045

									200.9732449102315,393.2795698924731

									211.21655613778938,387.36559139784947

									221.45986736534726,393.2795698924731

								"
								onclick="makeCellMove(8, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									88.29682140709487,209.94623655913978

									78.05351017953699,215.86021505376345

									67.81019895197912,209.94623655913978

									67.8101989519791,198.11827956989248

									78.05351017953699,192.2043010752688

									88.29682140709487,198.11827956989248

								"
								onclick="makeCellMove(9, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									98.54013263465275,227.68817204301075

									88.29682140709487,233.60215053763443

									78.053510179537,227.68817204301075

									78.05351017953699,215.86021505376345

									88.29682140709487,209.94623655913978

									98.54013263465275,215.86021505376345

								"
								onclick="makeCellMove(9, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									108.78344386221062,245.4301075268817

									98.54013263465274,251.34408602150535

									88.29682140709485,245.4301075268817

									88.29682140709485,233.6021505376344

									98.54013263465274,227.68817204301075

									108.78344386221062,233.6021505376344

								"
								onclick="makeCellMove(9, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									119.0267550897685,263.1720430107527

									108.78344386221062,269.0860215053763

									98.54013263465274,263.1720430107527

									98.54013263465274,251.34408602150538

									108.78344386221062,245.43010752688173

									119.0267550897685,251.34408602150538

								"
								onclick="makeCellMove(9, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,280.9139784946237

									119.0267550897685,286.8279569892473

									108.78344386221062,280.9139784946237

									108.78344386221062,269.0860215053763

									119.0267550897685,263.1720430107527

									129.27006631732638,269.0860215053763

								"
								onclick="makeCellMove(9, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,298.65591397849465

									129.27006631732638,304.5698924731183

									119.0267550897685,298.65591397849465

									119.0267550897685,286.8279569892473

									129.27006631732638,280.9139784946237

									139.51337754488426,286.8279569892473

								"
								onclick="makeCellMove(9, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,316.3978494623656

									139.51337754488424,322.31182795698925

									129.27006631732635,316.3978494623656

									129.27006631732635,304.5698924731183

									139.51337754488424,298.65591397849465

									149.75668877244212,304.5698924731183

								"
								onclick="makeCellMove(9, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,334.1397849462366

									149.75668877244212,340.0537634408602

									139.51337754488424,334.1397849462366

									139.51337754488424,322.31182795698925

									149.75668877244212,316.3978494623656

									160.0,322.31182795698925

								"
								onclick="makeCellMove(9, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,351.8817204301076

									160.0,357.7956989247312

									149.75668877244212,351.8817204301076

									149.75668877244212,340.0537634408602

									160.0,334.1397849462366

									170.24331122755788,340.0537634408602

								"
								onclick="makeCellMove(9, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,369.62365591397855

									170.24331122755788,375.5376344086022

									160.0,369.62365591397855

									160.0,357.7956989247312

									170.24331122755788,351.8817204301076

									180.48662245511576,357.7956989247312

								"
								onclick="makeCellMove(9, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,387.36559139784947

									180.48662245511576,393.2795698924731

									170.24331122755788,387.36559139784947

									170.24331122755788,375.5376344086021

									180.48662245511576,369.6236559139785

									190.72993368267365,375.5376344086021

								"
								onclick="makeCellMove(9, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,405.10752688172045

									190.72993368267362,411.02150537634407

									180.48662245511574,405.10752688172045

									180.48662245511574,393.2795698924731

									190.72993368267362,387.36559139784947

									200.9732449102315,393.2795698924731

								"
								onclick="makeCellMove(9, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									211.21655613778938,422.8494623655914

									200.9732449102315,428.76344086021504

									190.72993368267362,422.8494623655914

									190.72993368267362,411.02150537634407

									200.9732449102315,405.10752688172045

									211.21655613778938,411.02150537634407

								"
								onclick="makeCellMove(9, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									78.05351017953699,227.68817204301075

									67.8101989519791,233.60215053763443

									57.56688772442123,227.68817204301075

									57.56688772442122,215.86021505376345

									67.8101989519791,209.94623655913978

									78.05351017953699,215.86021505376345

								"
								onclick="makeCellMove(10, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									88.29682140709487,245.43010752688173

									78.05351017953699,251.3440860215054

									67.81019895197912,245.43010752688173

									67.8101989519791,233.60215053763443

									78.05351017953699,227.68817204301075

									88.29682140709487,233.60215053763443

								"
								onclick="makeCellMove(10, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									98.54013263465275,263.1720430107527

									88.29682140709487,269.0860215053763

									78.053510179537,263.1720430107527

									78.05351017953699,251.34408602150538

									88.29682140709487,245.43010752688173

									98.54013263465275,251.34408602150538

								"
								onclick="makeCellMove(10, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									108.78344386221062,280.9139784946237

									98.54013263465274,286.8279569892473

									88.29682140709485,280.9139784946237

									88.29682140709485,269.0860215053763

									98.54013263465274,263.1720430107527

									108.78344386221062,269.0860215053763

								"
								onclick="makeCellMove(10, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									119.0267550897685,298.65591397849465

									108.78344386221062,304.5698924731183

									98.54013263465274,298.65591397849465

									98.54013263465274,286.8279569892473

									108.78344386221062,280.9139784946237

									119.0267550897685,286.8279569892473

								"
								onclick="makeCellMove(10, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,316.3978494623656

									119.0267550897685,322.31182795698925

									108.78344386221062,316.3978494623656

									108.78344386221062,304.5698924731183

									119.0267550897685,298.65591397849465

									129.27006631732638,304.5698924731183

								"
								onclick="makeCellMove(10, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,334.1397849462366

									129.27006631732638,340.0537634408602

									119.0267550897685,334.1397849462366

									119.0267550897685,322.31182795698925

									129.27006631732638,316.3978494623656

									139.51337754488426,322.31182795698925

								"
								onclick="makeCellMove(10, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,351.8817204301076

									139.51337754488424,357.7956989247312

									129.27006631732635,351.8817204301076

									129.27006631732635,340.0537634408602

									139.51337754488424,334.1397849462366

									149.75668877244212,340.0537634408602

								"
								onclick="makeCellMove(10, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,369.62365591397855

									149.75668877244212,375.5376344086022

									139.51337754488424,369.62365591397855

									139.51337754488424,357.7956989247312

									149.75668877244212,351.8817204301076

									160.0,357.7956989247312

								"
								onclick="makeCellMove(10, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,387.3655913978495

									160.0,393.27956989247315

									149.75668877244212,387.3655913978495

									149.75668877244212,375.5376344086022

									160.0,369.62365591397855

									170.24331122755788,375.5376344086022

								"
								onclick="makeCellMove(10, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,405.1075268817205

									170.24331122755788,411.0215053763441

									160.0,405.1075268817205

									160.0,393.27956989247315

									170.24331122755788,387.3655913978495

									180.48662245511576,393.27956989247315

								"
								onclick="makeCellMove(10, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,422.8494623655914

									180.48662245511576,428.76344086021504

									170.24331122755788,422.8494623655914

									170.24331122755788,411.02150537634407

									180.48662245511576,405.10752688172045

									190.72993368267365,411.02150537634407

								"
								onclick="makeCellMove(10, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									200.9732449102315,440.5913978494624

									190.72993368267362,446.505376344086

									180.48662245511574,440.5913978494624

									180.48662245511574,428.76344086021504

									190.72993368267362,422.8494623655914

									200.9732449102315,428.76344086021504

								"
								onclick="makeCellMove(10, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									67.81019895197912,245.43010752688173

									57.566887724421235,251.3440860215054

									47.32357649686336,245.43010752688173

									47.32357649686335,233.60215053763443

									57.566887724421235,227.68817204301075

									67.81019895197912,233.60215053763443

								"
								onclick="makeCellMove(11, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									78.05351017953699,263.1720430107527

									67.8101989519791,269.0860215053763

									57.56688772442123,263.1720430107527

									57.56688772442122,251.34408602150538

									67.8101989519791,245.43010752688173

									78.05351017953699,251.34408602150538

								"
								onclick="makeCellMove(11, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									88.29682140709487,280.9139784946237

									78.05351017953699,286.8279569892473

									67.81019895197912,280.9139784946237

									67.8101989519791,269.0860215053763

									78.05351017953699,263.1720430107527

									88.29682140709487,269.0860215053763

								"
								onclick="makeCellMove(11, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									98.54013263465275,298.65591397849465

									88.29682140709487,304.5698924731183

									78.053510179537,298.65591397849465

									78.05351017953699,286.8279569892473

									88.29682140709487,280.9139784946237

									98.54013263465275,286.8279569892473

								"
								onclick="makeCellMove(11, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									108.78344386221062,316.3978494623656

									98.54013263465274,322.31182795698925

									88.29682140709485,316.3978494623656

									88.29682140709485,304.5698924731183

									98.54013263465274,298.65591397849465

									108.78344386221062,304.5698924731183

								"
								onclick="makeCellMove(11, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									119.0267550897685,334.1397849462366

									108.78344386221062,340.0537634408602

									98.54013263465274,334.1397849462366

									98.54013263465274,322.31182795698925

									108.78344386221062,316.3978494623656

									119.0267550897685,322.31182795698925

								"
								onclick="makeCellMove(11, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,351.8817204301076

									119.0267550897685,357.7956989247312

									108.78344386221062,351.8817204301076

									108.78344386221062,340.0537634408602

									119.0267550897685,334.1397849462366

									129.27006631732638,340.0537634408602

								"
								onclick="makeCellMove(11, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,369.62365591397855

									129.27006631732638,375.5376344086022

									119.0267550897685,369.62365591397855

									119.0267550897685,357.7956989247312

									129.27006631732638,351.8817204301076

									139.51337754488426,357.7956989247312

								"
								onclick="makeCellMove(11, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,387.3655913978495

									139.51337754488424,393.27956989247315

									129.27006631732635,387.3655913978495

									129.27006631732635,375.5376344086022

									139.51337754488424,369.62365591397855

									149.75668877244212,375.5376344086022

								"
								onclick="makeCellMove(11, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,405.1075268817205

									149.75668877244212,411.0215053763441

									139.51337754488424,405.1075268817205

									139.51337754488424,393.27956989247315

									149.75668877244212,387.3655913978495

									160.0,393.27956989247315

								"
								onclick="makeCellMove(11, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,422.8494623655915

									160.0,428.7634408602151

									149.75668877244212,422.8494623655915

									149.75668877244212,411.0215053763441

									160.0,405.1075268817205

									170.24331122755788,411.0215053763441

								"
								onclick="makeCellMove(11, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,440.59139784946245

									170.24331122755788,446.5053763440861

									160.0,440.59139784946245

									160.0,428.7634408602151

									170.24331122755788,422.8494623655915

									180.48662245511576,428.7634408602151

								"
								onclick="makeCellMove(11, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									190.72993368267365,458.33333333333337

									180.48662245511576,464.247311827957

									170.24331122755788,458.33333333333337

									170.24331122755788,446.505376344086

									180.48662245511576,440.5913978494624

									190.72993368267365,446.505376344086

								"
								onclick="makeCellMove(11, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									57.56688772442123,263.1720430107527

									47.32357649686335,269.0860215053763

									37.08026526930548,263.1720430107527

									37.08026526930547,251.34408602150538

									47.32357649686335,245.43010752688173

									57.566887724421235,251.34408602150538

								"
								onclick="makeCellMove(12, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									67.81019895197912,280.9139784946237

									57.566887724421235,286.8279569892473

									47.32357649686336,280.9139784946237

									47.32357649686335,269.0860215053763

									57.566887724421235,263.1720430107527

									67.81019895197912,269.0860215053763

								"
								onclick="makeCellMove(12, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									78.05351017953699,298.65591397849465

									67.8101989519791,304.5698924731183

									57.56688772442123,298.65591397849465

									57.56688772442122,286.8279569892473

									67.8101989519791,280.9139784946237

									78.05351017953699,286.8279569892473

								"
								onclick="makeCellMove(12, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									88.29682140709487,316.3978494623656

									78.05351017953699,322.31182795698925

									67.81019895197912,316.3978494623656

									67.8101989519791,304.5698924731183

									78.05351017953699,298.65591397849465

									88.29682140709487,304.5698924731183

								"
								onclick="makeCellMove(12, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									98.54013263465275,334.1397849462366

									88.29682140709487,340.0537634408602

									78.053510179537,334.1397849462366

									78.05351017953699,322.31182795698925

									88.29682140709487,316.3978494623656

									98.54013263465275,322.31182795698925

								"
								onclick="makeCellMove(12, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									108.78344386221062,351.8817204301075

									98.54013263465274,357.79569892473114

									88.29682140709485,351.8817204301075

									88.29682140709485,340.05376344086017

									98.54013263465274,334.13978494623655

									108.78344386221062,340.05376344086017

								"
								onclick="makeCellMove(12, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									119.0267550897685,369.6236559139785

									108.78344386221062,375.5376344086021

									98.54013263465274,369.6236559139785

									98.54013263465274,357.79569892473114

									108.78344386221062,351.8817204301075

									119.0267550897685,357.79569892473114

								"
								onclick="makeCellMove(12, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,387.36559139784947

									119.0267550897685,393.2795698924731

									108.78344386221062,387.36559139784947

									108.78344386221062,375.5376344086021

									119.0267550897685,369.6236559139785

									129.27006631732638,375.5376344086021

								"
								onclick="makeCellMove(12, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,405.10752688172045

									129.27006631732638,411.02150537634407

									119.0267550897685,405.10752688172045

									119.0267550897685,393.2795698924731

									129.27006631732638,387.36559139784947

									139.51337754488426,393.2795698924731

								"
								onclick="makeCellMove(12, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,422.8494623655914

									139.51337754488424,428.76344086021504

									129.27006631732635,422.8494623655914

									129.27006631732635,411.02150537634407

									139.51337754488424,405.10752688172045

									149.75668877244212,411.02150537634407

								"
								onclick="makeCellMove(12, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,440.5913978494624

									149.75668877244212,446.505376344086

									139.51337754488424,440.5913978494624

									139.51337754488424,428.76344086021504

									149.75668877244212,422.8494623655914

									160.0,428.76344086021504

								"
								onclick="makeCellMove(12, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,458.3333333333333

									160.0,464.24731182795693

									149.75668877244212,458.3333333333333

									149.75668877244212,446.50537634408596

									160.0,440.59139784946234

									170.24331122755788,446.50537634408596

								"
								onclick="makeCellMove(12, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									180.48662245511576,476.0752688172043

									170.24331122755788,481.9892473118279

									160.0,476.0752688172043

									160.0,464.24731182795693

									170.24331122755788,458.3333333333333

									180.48662245511576,464.24731182795693

								"
								onclick="makeCellMove(12, 13, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									47.323576496863346,280.9139784946237

									37.08026526930547,286.8279569892473

									26.836954041747596,280.9139784946237

									26.836954041747592,269.0860215053763

									37.08026526930547,263.1720430107527

									47.32357649686335,269.0860215053763

								"
								onclick="makeCellMove(13, 1, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									57.56688772442123,298.65591397849465

									47.32357649686335,304.5698924731183

									37.08026526930548,298.65591397849465

									37.08026526930547,286.8279569892473

									47.32357649686335,280.9139784946237

									57.566887724421235,286.8279569892473

								"
								onclick="makeCellMove(13, 2, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									67.81019895197912,316.3978494623656

									57.566887724421235,322.31182795698925

									47.32357649686336,316.3978494623656

									47.32357649686335,304.5698924731183

									57.566887724421235,298.65591397849465

									67.81019895197912,304.5698924731183

								"
								onclick="makeCellMove(13, 3, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									78.05351017953699,334.1397849462366

									67.8101989519791,340.0537634408602

									57.56688772442123,334.1397849462366

									57.56688772442122,322.31182795698925

									67.8101989519791,316.3978494623656

									78.05351017953699,322.31182795698925

								"
								onclick="makeCellMove(13, 4, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									88.29682140709487,351.8817204301076

									78.05351017953699,357.7956989247312

									67.81019895197912,351.8817204301076

									67.8101989519791,340.0537634408602

									78.05351017953699,334.1397849462366

									88.29682140709487,340.0537634408602

								"
								onclick="makeCellMove(13, 5, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									98.54013263465275,369.62365591397855

									88.29682140709487,375.5376344086022

									78.053510179537,369.62365591397855

									78.05351017953699,357.7956989247312

									88.29682140709487,351.8817204301076

									98.54013263465275,357.7956989247312

								"
								onclick="makeCellMove(13, 6, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									108.78344386221062,387.36559139784947

									98.54013263465274,393.2795698924731

									88.29682140709485,387.36559139784947

									88.29682140709485,375.5376344086021

									98.54013263465274,369.6236559139785

									108.78344386221062,375.5376344086021

								"
								onclick="makeCellMove(13, 7, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									119.0267550897685,405.10752688172045

									108.78344386221062,411.02150537634407

									98.54013263465274,405.10752688172045

									98.54013263465274,393.2795698924731

									108.78344386221062,387.36559139784947

									119.0267550897685,393.2795698924731

								"
								onclick="makeCellMove(13, 8, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									129.27006631732638,422.8494623655914

									119.0267550897685,428.76344086021504

									108.78344386221062,422.8494623655914

									108.78344386221062,411.02150537634407

									119.0267550897685,405.10752688172045

									129.27006631732638,411.02150537634407

								"
								onclick="makeCellMove(13, 9, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									139.51337754488426,440.5913978494624

									129.27006631732638,446.505376344086

									119.0267550897685,440.5913978494624

									119.0267550897685,428.76344086021504

									129.27006631732638,422.8494623655914

									139.51337754488426,428.76344086021504

								"
								onclick="makeCellMove(13, 10, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									149.75668877244212,458.33333333333337

									139.51337754488424,464.247311827957

									129.27006631732635,458.33333333333337

									129.27006631732635,446.505376344086

									139.51337754488424,440.5913978494624

									149.75668877244212,446.505376344086

								"
								onclick="makeCellMove(13, 11, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									160.0,476.07526881720435

									149.75668877244212,481.98924731182797

									139.51337754488424,476.07526881720435

									139.51337754488424,464.247311827957

									149.75668877244212,458.33333333333337

									160.0,464.247311827957

								"
								onclick="makeCellMove(13, 12, player);"
								style="stroke: black; fill: white;"
							/>

							<polygon
								points="

									170.24331122755788,493.81720430107526

									160.0,499.7311827956989

									149.75668877244212,493.81720430107526

									149.75668877244212,481.9892473118279

									160.0,476.0752688172043

									170.24331122755788,481.9892473118279

								"
								onclick="makeCellMove(13, 13, player);"
								style="stroke: black; fill: white;"
							/>

					</svg>
  				</td>
  				<td style="text-align: center;">
  				<table>
  				
  					<tr>
  						<td style="text-align: right;">
			  				
								<span style="color: white;">Opposite</span>
							
		  				</td>
			  			<td colspan="2">
							<div style="position: relative;">
								<div id="loader"></div>
								<img id="opponentFlipCardImg" src="/static/cardstux/blank.png" width="54" height="72">
								<img id="opponentPlayCardImg" src="/static/cardstux/blank.png" width="54" height="72"><br />
							</div>
	  			  		</td>
  				  		<td></td>
  				  	</tr>
					<tr>
  				  		<td style="text-align: right;">	
		  				
							<span style="color: white;">Same</span>
						
			  			</td>
			  			<td colspan="2">
							<img id="playerFlipCardImg" src="/static/cardstux/blank.png" width="54" height="72">
							<img id="playerPlayCardImg" src="/static/cardstux/blank.png" width="54" height="72"><br />
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
    _ = pagePartSVGLabels2
}
