var colors = ["black", "white"];
var color = colors[Math.floor(Math.random() * colors.length)];
var humansTurn = color === "white";

function onDrop(source, target, piece, newPos, oldPos, orientation) {
  if (piece[0] !== orientation[0]) return "snapback"; // only move your own pieces
  if (!humansTurn) return "snapback"; // only move if it's your turn
  if (target === "offboard") return "snapback"; // only move pieces onto the board
  humansTurn = false;
  postMoveToThink(source + target);
}

var cfg = {
  draggable: true,
  dropOffBoard: "snapback",
  position: "start",
  onDrop: onDrop,
  orientation: color
};

var board1 = ChessBoard("board1", cfg);

function postMoveToThink(move) {
  console.log("posting move", move);
  $.ajax({
    url: "/think",
    data: { move },
    method: "post",
    success: onSucces
  });
}

function restartGame() {
  $.ajax({ url: "/restart", method: "post" });
}

function onSucces(response) {
  var fen = JSON.parse(response);
  console.log(fen);
  board1.position(fen, true);
  humansTurn = true;
}

function beginGame() {
  $.ajax({
    url: "/begin",
    method: "post",
    success: onSucces
  });
}

restartGame();
if (!humansTurn) beginGame();
