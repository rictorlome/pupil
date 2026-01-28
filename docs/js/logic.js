var wasmReady = false;
var board1 = null;
var colors = ["black", "white"];
var color = colors[Math.floor(Math.random() * colors.length)];
var humansTurn = color === "white";

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

  board1.position(result.fen, true);
  setStatus("Engine thinking...");

  // Let the engine respond (use setTimeout to not block UI)
  setTimeout(() => {
    var engineResult = pupilGetEngineMove(5); // depth 5
    console.log("Engine played:", engineResult.move);
    board1.position(engineResult.fen, true);
    humansTurn = true;
    setStatus("Your turn!");
  }, 50);
}

function restartGame() {
  if (!wasmReady) return;

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
      var engineResult = pupilGetEngineMove(5);
      console.log("Engine played:", engineResult.move);
      board1.position(engineResult.fen, true);
      humansTurn = true;
      setStatus("Your turn!");
    }, 50);
  } else {
    setStatus("Your turn!");
  }
}
