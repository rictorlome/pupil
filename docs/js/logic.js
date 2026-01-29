var wasmReady = false;
var board1 = null;
var colors = ["black", "white"];
var color = colors[Math.floor(Math.random() * colors.length)];
var humansTurn = color === "white";
var gameOver = false;

// Load WASM
const go = new Go();
WebAssembly.instantiateStreaming(fetch("pupil.wasm"), go.importObject).then((result) => {
  go.run(result.instance);
  wasmReady = true;
  initGame();
}).catch((err) => {
  console.error("Failed to load WASM:", err);
  setStatus("Failed to load engine: " + err.message);
});

function setStatus(msg) {
  document.getElementById("status").textContent = msg;
}

function initGame() {
  setStatus("Engine loaded! Starting game...");

  var cfg = {
    draggable: true,
    dropOffBoard: "snapback",
    position: "start",
    onDrop: onDrop,
    orientation: color
  };

  board1 = ChessBoard("board1", cfg);
  restartGame();
}

function onDrop(source, target, piece, newPos, oldPos, orientation) {
  if (!wasmReady) return "snapback";
  if (gameOver) return "snapback";
  if (piece[0] !== orientation[0]) return "snapback";
  if (!humansTurn) return "snapback";
  if (target === "offboard") return "snapback";

  humansTurn = false;
  setStatus("Thinking...");

  // Use setTimeout to let the UI update before the engine thinks
  setTimeout(() => {
    makeMove(source + target);
  }, 10);
}

function checkGameStatus() {
  var status = pupilGetGameStatus();

  if (status.status === "checkmate") {
    gameOver = true;
    if (status.winner === color) {
      setStatus("Checkmate! You win!");
    } else {
      setStatus("Checkmate! Engine wins!");
    }
    return true;
  }

  if (status.status === "stalemate") {
    gameOver = true;
    setStatus("Stalemate! It's a draw.");
    return true;
  }

  return false;
}

function makeMove(move) {
  console.log("Making move:", move);

  // Make the human's move
  var result = pupilMakeMove(move);

  if (result.error) {
    console.log("Illegal move:", result.error);
    board1.position(result.fen || pupilGetFen(), true);
    humansTurn = true;
    setStatus("Illegal move! Your turn.");
    return;
  }

  console.log("After human move - FEN:", result.fen);
  board1.position(result.fen, true);

  // Check if game ended after human's move
  if (checkGameStatus()) {
    console.log("Game over after human move");
    return;
  }

  setStatus("Engine thinking...");

  // Let the engine respond (use setTimeout to not block UI)
  setTimeout(() => {
    var engineResult = pupilGetEngineMove();
    console.log("Engine played:", engineResult.move);
    console.log("After engine move - FEN:", engineResult.fen);
    board1.position(engineResult.fen, true);

    // Check if game ended after engine's move
    if (checkGameStatus()) {
      console.log("Game over after engine move");
      return;
    }

    humansTurn = true;
    var status = pupilGetGameStatus();
    if (status.inCheck) {
      setStatus("Check! Your turn.");
    } else {
      setStatus("Your turn!");
    }
  }, 50);
}

function restartGame() {
  if (!wasmReady) return;

  gameOver = false;
  color = colors[Math.floor(Math.random() * colors.length)];
  humansTurn = color === "white";

  var fen = pupilNewGame();
  console.log("New game started:", fen);

  if (board1) {
    board1.orientation(color);
    board1.position("start", false);
  }

  if (!humansTurn) {
    setStatus("Engine thinking...");
    setTimeout(() => {
      var engineResult = pupilGetEngineMove();
      console.log("Engine played:", engineResult.move);
      board1.position(engineResult.fen, true);
      humansTurn = true;
      setStatus("Your turn!");
    }, 50);
  } else {
    setStatus("Your turn!");
  }
}
