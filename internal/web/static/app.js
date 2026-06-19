const COLORS = ["R", "G", "B", "Y", "O", "P", "C", "M"];

const colorMap = {
  R: "#e53935",
  G: "#43a047",
  B: "#1e88e5",
  Y: "#fdd835",
  O: "#fb8c00",
  P: "#8e24aa",
  C: "#00acc1",
  M: "#d81b60",
};

const DEFAULT_PUZZLE = `5 5
R 0 2 4 2
G 2 0 2 4
X 2 2`;

const state = {
  rows: 5,
  cols: 5,
  tool: "dot",
  activeColor: "R",
  endpoints: {},
  bridges: new Set(),
  walls: new Set(),
};

function cellKey(r, c) {
  return `${r},${c}`;
}

function parseKey(key) {
  const [r, c] = key.split(",").map(Number);
  return { r, c };
}

function emptyEndpoints() {
  const ep = {};
  for (const ch of COLORS) ep[ch] = [];
  return ep;
}

function resetState(rows, cols, keepContent = false) {
  const old = keepContent
    ? { endpoints: state.endpoints, bridges: state.bridges, walls: state.walls }
    : null;
  state.rows = rows;
  state.cols = cols;
  state.endpoints = emptyEndpoints();
  state.bridges = new Set();
  state.walls = new Set();
  if (keepContent && old) {
    for (const ch of COLORS) {
      state.endpoints[ch] = (old.endpoints[ch] || []).filter(
        (p) => p.r < rows && p.c < cols
      );
    }
    for (const key of old.bridges) {
      const { r, c } = parseKey(key);
      if (r < rows && c < cols) state.bridges.add(key);
    }
    for (const key of old.walls) {
      const { r, c } = parseKey(key);
      if (r < rows && c < cols) state.walls.add(key);
    }
  }
}

function endpointAt(r, c) {
  for (const ch of COLORS) {
    const idx = state.endpoints[ch].findIndex((p) => p.r === r && p.c === c);
    if (idx >= 0) return { color: ch, index: idx };
  }
  return null;
}

function removeAllAt(r, c) {
  const key = cellKey(r, c);
  state.bridges.delete(key);
  state.walls.delete(key);
  for (const ch of COLORS) {
    state.endpoints[ch] = state.endpoints[ch].filter(
      (p) => !(p.r === r && p.c === c)
    );
  }
}

function applyDot(r, c) {
  const key = cellKey(r, c);
  if (state.walls.has(key) || state.bridges.has(key)) return;
  const existing = endpointAt(r, c);
  if (existing) {
    if (existing.color === state.activeColor) {
      state.endpoints[existing.color].splice(existing.index, 1);
    }
    return;
  }
  const list = state.endpoints[state.activeColor];
  if (list.length >= 2) list.shift();
  list.push({ r, c });
}

function applyBridge(r, c) {
  const key = cellKey(r, c);
  if (state.walls.has(key) || endpointAt(r, c)) return;
  if (state.bridges.has(key)) state.bridges.delete(key);
  else state.bridges.add(key);
}

function applyWall(r, c) {
  const key = cellKey(r, c);
  removeAllAt(r, c);
  if (state.walls.has(key)) state.walls.delete(key);
  else state.walls.add(key);
}

function handleCellClick(r, c) {
  switch (state.tool) {
    case "dot":
      applyDot(r, c);
      break;
    case "bridge":
      applyBridge(r, c);
      break;
    case "wall":
      applyWall(r, c);
      break;
    case "erase":
      removeAllAt(r, c);
      break;
  }
  renderEditor();
  renderPalette();
}

function renderPalette() {
  const el = document.getElementById("palette");
  el.innerHTML = "";
  for (const ch of COLORS) {
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "swatch" + (state.activeColor === ch ? " active" : "");
    btn.style.background = colorMap[ch];
    btn.title = `Color ${ch}`;
    const count = state.endpoints[ch].length;
    btn.innerHTML = `<span class="count">${count}/2</span>`;
    btn.addEventListener("click", () => {
      state.activeColor = ch;
      state.tool = "dot";
      document.querySelectorAll(".tool").forEach((t) => {
        t.classList.toggle("active", t.dataset.tool === "dot");
      });
      renderPalette();
    });
    el.appendChild(btn);
  }
}

function renderEditor() {
  const grid = document.getElementById("editor-grid");
  grid.innerHTML = "";
  grid.style.gridTemplateColumns = `repeat(${state.cols}, 44px)`;
  for (let r = 0; r < state.rows; r++) {
    for (let c = 0; c < state.cols; c++) {
      const cell = document.createElement("div");
      cell.className = "editor-cell";
      cell.dataset.r = r;
      cell.dataset.c = c;
      const key = cellKey(r, c);
      if (state.walls.has(key)) cell.classList.add("wall");
      if (state.bridges.has(key)) cell.classList.add("bridge");
      const ep = endpointAt(r, c);
      if (ep) {
        const dot = document.createElement("span");
        dot.className = "dot";
        dot.style.background = colorMap[ep.color];
        cell.appendChild(dot);
      }
      cell.addEventListener("click", () => handleCellClick(r, c));
      grid.appendChild(cell);
    }
  }
}

function toPuzzleText() {
  const lines = [`${state.rows} ${state.cols}`];
  for (const ch of COLORS) {
    const pts = state.endpoints[ch];
    if (pts.length === 2) {
      lines.push(`${ch} ${pts[0].r} ${pts[0].c} ${pts[1].r} ${pts[1].c}`);
    }
  }
  const bridgeLines = [...state.bridges]
    .map(parseKey)
    .sort((a, b) => a.r - b.r || a.c - b.c)
    .map(({ r, c }) => `X ${r} ${c}`);
  const wallLines = [...state.walls]
    .map(parseKey)
    .sort((a, b) => a.r - b.r || a.c - b.c)
    .map(({ r, c }) => `W ${r} ${c}`);
  lines.push(...bridgeLines, ...wallLines);
  return lines.join("\n");
}

function parsePuzzleText(text) {
  const lines = text
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l && !l.startsWith("#"));
  if (lines.length === 0) throw new Error("empty puzzle");
  const head = lines[0].split(/\s+/);
  const rows = parseInt(head[0], 10);
  const cols = parseInt(head[1], 10);
  if (!rows || !cols) throw new Error("first line must be: rows cols");

  resetState(rows, cols);
  for (let i = 1; i < lines.length; i++) {
    const p = lines[i].split(/\s+/);
    if (p[0] === "X" || p[0] === "x") {
      state.bridges.add(cellKey(+p[1], +p[2]));
    } else if (p[0] === "W" || p[0] === "w") {
      state.walls.add(cellKey(+p[1], +p[2]));
    } else if (p.length >= 5) {
      const ch = p[0][0];
      if (!state.endpoints[ch]) state.endpoints[ch] = [];
      state.endpoints[ch].push({ r: +p[1], c: +p[2] }, { r: +p[3], c: +p[4] });
    }
  }
}

function bridgeSetFromGrid(gridJSON) {
  const set = new Set();
  if (!gridJSON || !gridJSON.grid) return set;
  for (let r = 0; r < gridJSON.rows; r++) {
    for (let c = 0; c < gridJSON.cols; c++) {
      const ch = gridJSON.grid[r][c];
      if (ch === "x" || ch === "+") set.add(cellKey(r, c));
    }
  }
  return set;
}

function endpointSetFromGrid(gridJSON) {
  const set = new Set();
  if (!gridJSON || !gridJSON.grid) return set;
  for (let r = 0; r < gridJSON.rows; r++) {
    for (let c = 0; c < gridJSON.cols; c++) {
      const ch = gridJSON.grid[r][c];
      if (isColor(ch)) set.add(cellKey(r, c));
    }
  }
  return set;
}

function isBridgeAt(r, c, bridgeSet) {
  if (bridgeSet) return bridgeSet.has(cellKey(r, c));
  return state.bridges.has(cellKey(r, c));
}

function endpointSetFromState() {
  const set = new Set();
  for (const ch of COLORS) {
    for (const p of state.endpoints[ch] || []) {
      set.add(cellKey(p.r, p.c));
    }
  }
  return set;
}

function gridCells(gridJSON) {
  return gridJSON.grid.map((row) => [...row]);
}

function isColor(ch) {
  return !!colorMap[ch];
}

function axisBetween(r1, c1, r2, c2) {
  if (r1 === r2) return "h";
  if (c1 === c2) return "v";
  return null;
}

function horizColor(cells, r, c, cols) {
  const left = c > 0 ? cells[r][c - 1] : ".";
  const right = c < cols - 1 ? cells[r][c + 1] : ".";
  if (isColor(left) && left === right) return left;
  if (isColor(left)) return left;
  if (isColor(right)) return right;
  return null;
}

function vertColor(cells, r, c, rows) {
  const up = r > 0 ? cells[r - 1][c] : ".";
  const down = r < rows - 1 ? cells[r + 1][c] : ".";
  if (isColor(up) && up === down) return up;
  if (isColor(up)) return up;
  if (isColor(down)) return down;
  return null;
}

function trackColor(cells, r, c, rows, cols, axis, bridgeSet) {
  const ch = cells[r][c];
  if (ch === "#" || ch === "." || ch === "x") return null;
  if (ch === "+") {
    return axis === "h" ? horizColor(cells, r, c, cols) : vertColor(cells, r, c, rows);
  }
  if (isColor(ch)) {
    if (isBridgeAt(r, c, bridgeSet)) {
      if (axis === "h") {
        const left = c > 0 ? cells[r][c - 1] : ".";
        const right = c < cols - 1 ? cells[r][c + 1] : ".";
        if (left === ch || right === ch) return ch;
        return null;
      }
      const up = r > 0 ? cells[r - 1][c] : ".";
      const down = r < rows - 1 ? cells[r + 1][c] : ".";
      if (up === ch || down === ch) return ch;
      return null;
    }
    return ch;
  }
  return null;
}

function pipesConnect(cells, r1, c1, r2, c2, rows, cols, bridgeSet) {
  const axis = axisBetween(r1, c1, r2, c2);
  if (!axis) return false;
  const a = trackColor(cells, r1, c1, rows, cols, axis, bridgeSet);
  const b = trackColor(cells, r2, c2, rows, cols, axis, bridgeSet);
  return a && a === b;
}

function edgeKey(r1, c1, r2, c2) {
  const a = `${r1},${c1}`;
  const b = `${r2},${c2}`;
  return a < b ? `${a}|${b}` : `${b}|${a}`;
}

function hasPathEdge(pathEdges, r1, c1, r2, c2) {
  return pathEdges.has(edgeKey(r1, c1, r2, c2));
}

function cellHasColor(cells, r, c, color, rows, cols) {
  const ch = cells[r][c];
  if (ch === color) return true;
  if (ch === "+") {
    return (
      horizColor(cells, r, c, cols) === color ||
      vertColor(cells, r, c, rows) === color
    );
  }
  return false;
}

function pathNeighbors(r, c, color, cells, rows, cols, bridgeSet) {
  const out = [];
  const dirs = [
    [-1, 0],
    [0, 1],
    [1, 0],
    [0, -1],
  ];
  for (const [dr, dc] of dirs) {
    const nr = r + dr;
    const nc = c + dc;
    if (nr < 0 || nc < 0 || nr >= rows || nc >= cols) continue;
    if (!cellHasColor(cells, nr, nc, color, rows, cols)) continue;
    if (pipesConnect(cells, r, c, nr, nc, rows, cols, bridgeSet)) {
      out.push({ r: nr, c: nc });
    }
  }
  return out;
}

function collectColorCells(color, cells, rows, cols) {
  const list = [];
  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      if (cellHasColor(cells, r, c, color, rows, cols)) {
        list.push({ r, c });
      }
    }
  }
  return list;
}

function endpointsForColor(color, inputGrid, epSet) {
  const eps = [];
  if (inputGrid && inputGrid.grid) {
    for (let r = 0; r < inputGrid.rows; r++) {
      for (let c = 0; c < inputGrid.cols; c++) {
        if (inputGrid.grid[r][c] === color) eps.push({ r, c });
      }
    }
    return eps;
  }
  for (const key of epSet) {
    const { r, c } = parseKey(key);
    for (const ch of COLORS) {
      if (ch === color && state.endpoints[ch].some((p) => p.r === r && p.c === c)) {
        eps.push({ r, c });
      }
    }
  }
  return eps;
}

function tracePathEdges(color, cells, rows, cols, bridgeSet, inputGrid, epSet) {
  const colorCells = collectColorCells(color, cells, rows, cols);
  if (colorCells.length <= 1) return new Set();

  const endpoints = endpointsForColor(color, inputGrid, epSet);
  if (endpoints.length !== 2) return new Set();

  const targetSize = colorCells.length;
  const endKey = cellKey(endpoints[1].r, endpoints[1].c);
  const pathEdges = new Set();

  function dfs(r, c, prev, visited, edges) {
    const key = cellKey(r, c);
    if (visited.has(key)) return false;
    visited.add(key);

    if (visited.size === targetSize) {
      if (key === endKey) return true;
      visited.delete(key);
      return false;
    }

    for (const n of pathNeighbors(r, c, color, cells, rows, cols, bridgeSet)) {
      if (prev && n.r === prev.r && n.c === prev.c) continue;
      if (dfs(n.r, n.c, { r, c }, visited, edges)) {
        edges.add(edgeKey(r, c, n.r, n.c));
        return true;
      }
    }

    visited.delete(key);
    return false;
  }

  if (dfs(endpoints[0].r, endpoints[0].c, null, new Set(), pathEdges)) {
    return pathEdges;
  }
  const reverseEdges = new Set();
  if (dfs(endpoints[1].r, endpoints[1].c, null, new Set(), reverseEdges)) {
    return reverseEdges;
  }
  return fallbackPathEdges(color, cells, rows, cols, bridgeSet, endpoints);
}

function fallbackPathEdges(color, cells, rows, cols, bridgeSet, endpoints) {
  const pathEdges = new Set();
  let cur = endpoints[0];
  let prev = null;
  const visited = new Set([cellKey(cur.r, cur.c)]);

  while (visited.size < collectColorCells(color, cells, rows, cols).length) {
    const nexts = pathNeighbors(cur.r, cur.c, color, cells, rows, cols, bridgeSet).filter(
      (n) =>
        (!prev || n.r !== prev.r || n.c !== prev.c) &&
        !visited.has(cellKey(n.r, n.c))
    );
    if (nexts.length === 0) break;
    const next = nexts[0];
    pathEdges.add(edgeKey(cur.r, cur.c, next.r, next.c));
    prev = cur;
    cur = next;
    visited.add(cellKey(cur.r, cur.c));
  }
  return pathEdges;
}

function buildPathEdgeSet(cells, rows, cols, bridgeSet, inputGrid, epSet, colors) {
  const pathEdges = new Set();
  const palette = colors && colors.length ? colors : COLORS;
  for (const color of palette) {
    if (!isColor(color)) continue;
    for (const key of tracePathEdges(color, cells, rows, cols, bridgeSet, inputGrid, epSet)) {
      pathEdges.add(key);
    }
  }
  return pathEdges;
}

function appendPipeArm(cell, dir, color, trackClass) {
  const arm = document.createElement("div");
  arm.className = "pipe-arm " + dir + (trackClass ? " " + trackClass : "");
  if (trackClass === "h-track") {
    cell.style.setProperty("--pipe-color-h", colorMap[color]);
  } else if (trackClass === "v-track") {
    cell.style.setProperty("--pipe-color-v", colorMap[color]);
  } else {
    cell.style.setProperty("--pipe-color", colorMap[color]);
  }
  cell.appendChild(arm);
}

function appendPipeCore(cell, color, isEndpoint) {
  const core = document.createElement("div");
  core.className = "pipe-core" + (isEndpoint ? " endpoint" : "");
  cell.style.setProperty("--pipe-color", colorMap[color]);
  cell.appendChild(core);
}

function renderPipeCell(cell, cells, r, c, rows, cols, color, epSet, pathEdges) {
  const dirs = [
    ["n", -1, 0],
    ["e", 0, 1],
    ["s", 1, 0],
    ["w", 0, -1],
  ];
  for (const [name, dr, dc] of dirs) {
    const nr = r + dr;
    const nc = c + dc;
    if (nr < 0 || nc < 0 || nr >= rows || nc >= cols) continue;
    if (hasPathEdge(pathEdges, r, c, nr, nc)) {
      appendPipeArm(cell, name, color);
    }
  }
  appendPipeCore(cell, color, epSet.has(cellKey(r, c)));
}

function renderBridgeCross(cell, cells, r, c, rows, cols, pathEdges) {
  cell.classList.add("pipe-cross");
  const hCol = horizColor(cells, r, c, cols);
  const vCol = vertColor(cells, r, c, rows);
  if (hCol) {
    cell.style.setProperty("--pipe-color-h", colorMap[hCol]);
    if (c > 0 && hasPathEdge(pathEdges, r, c, r, c - 1)) {
      appendPipeArm(cell, "w", hCol, "h-track");
    }
    if (c < cols - 1 && hasPathEdge(pathEdges, r, c, r, c + 1)) {
      appendPipeArm(cell, "e", hCol, "h-track");
    }
    const coreH = document.createElement("div");
    coreH.className = "pipe-core h-track";
    cell.appendChild(coreH);
  }
  if (vCol) {
    cell.style.setProperty("--pipe-color-v", colorMap[vCol]);
    if (r > 0 && hasPathEdge(pathEdges, r, c, r - 1, c)) {
      appendPipeArm(cell, "n", vCol, "v-track");
    }
    if (r < rows - 1 && hasPathEdge(pathEdges, r, c, r + 1, c)) {
      appendPipeArm(cell, "s", vCol, "v-track");
    }
    const coreV = document.createElement("div");
    coreV.className = "pipe-core v-track";
    cell.appendChild(coreV);
  }
}

function renderPipeGrid(el, gridJSON, inputGrid) {
  el.innerHTML = "";
  if (!gridJSON || !gridJSON.grid) return;
  const rows = gridJSON.rows;
  const cols = gridJSON.cols;
  const cells = gridCells(gridJSON);
  const bridgeSet = bridgeSetFromGrid(inputGrid || gridJSON);
  const epSet = inputGrid
    ? endpointSetFromGrid(inputGrid)
    : endpointSetFromState();
  const pathEdges = buildPathEdgeSet(
    cells,
    rows,
    cols,
    bridgeSet,
    inputGrid,
    epSet,
    gridJSON.colors
  );

  el.className = "grid preview pipes";
  el.style.gridTemplateColumns = `repeat(${cols}, 32px)`;

  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      const cell = document.createElement("div");
      cell.className = "pipe-cell";
      const ch = cells[r][c];

      if (ch === "#") {
        cell.classList.add("wall");
        el.appendChild(cell);
        continue;
      }
      if (ch === "x") {
        cell.classList.add("bridge-empty");
        el.appendChild(cell);
        continue;
      }
      if (ch === "+") {
        renderBridgeCross(cell, cells, r, c, rows, cols, pathEdges);
        el.appendChild(cell);
        continue;
      }
      if (ch === ".") {
        cell.classList.add("empty");
        el.appendChild(cell);
        continue;
      }
      if (isColor(ch)) {
        renderPipeCell(cell, cells, r, c, rows, cols, ch, epSet, pathEdges);
        el.appendChild(cell);
        continue;
      }
      cell.classList.add("empty");
      el.appendChild(cell);
    }
  }
}

function renderInputPreview(el, inputGrid) {
  el.innerHTML = "";
  if (!inputGrid || !inputGrid.grid) {
    renderInputPreviewFromState(el);
    return;
  }
  el.className = "grid preview pipes";
  el.style.gridTemplateColumns = `repeat(${inputGrid.cols}, 32px)`;

  for (let r = 0; r < inputGrid.rows; r++) {
    for (let c = 0; c < inputGrid.cols; c++) {
      const cell = document.createElement("div");
      cell.className = "pipe-cell";
      const ch = inputGrid.grid[r][c];

      if (ch === "#") {
        cell.classList.add("wall");
      } else if (ch === "x") {
        cell.classList.add("bridge-empty");
      } else {
        cell.classList.add("empty");
      }

      if (isColor(ch)) {
        const dot = document.createElement("div");
        dot.className = "pipe-core endpoint";
        dot.style.setProperty("--pipe-color", colorMap[ch]);
        dot.style.background = colorMap[ch];
        cell.appendChild(dot);
      }

      el.appendChild(cell);
    }
  }
}

function renderInputPreviewFromState(el) {
  el.innerHTML = "";
  el.className = "grid preview pipes";
  el.style.gridTemplateColumns = `repeat(${state.cols}, 32px)`;

  for (let r = 0; r < state.rows; r++) {
    for (let c = 0; c < state.cols; c++) {
      const cell = document.createElement("div");
      cell.className = "pipe-cell";
      const key = cellKey(r, c);

      if (state.walls.has(key)) {
        cell.classList.add("wall");
      } else if (state.bridges.has(key)) {
        cell.classList.add("bridge-empty");
      } else {
        cell.classList.add("empty");
      }

      const ep = endpointAt(r, c);
      if (ep) {
        const dot = document.createElement("div");
        dot.className = "pipe-core endpoint";
        dot.style.setProperty("--pipe-color", colorMap[ep.color]);
        dot.style.background = colorMap[ep.color];
        cell.appendChild(dot);
      }

      el.appendChild(cell);
    }
  }
}

function renderGrid(el, gridJSON, inputGrid) {
  renderPipeGrid(el, gridJSON, inputGrid);
}

async function copyText(text, statusEl) {
  try {
    await navigator.clipboard.writeText(text);
    if (statusEl) {
      statusEl.textContent = "Copied!";
      setTimeout(() => (statusEl.textContent = ""), 2000);
    }
  } catch {
    if (statusEl) statusEl.textContent = "Copy failed";
  }
}

function getPuzzleText() {
  const textPanel = document.getElementById("panel-text");
  if (textPanel.classList.contains("active")) {
    return document.getElementById("puzzle").value;
  }
  return toPuzzleText();
}

async function solve() {
  const puzzle = getPuzzleText();
  const status = document.getElementById("status");
  status.textContent = "Solving...";
  try {
    parsePuzzleText(puzzle);
    document.getElementById("rows").value = state.rows;
    document.getElementById("cols").value = state.cols;
  } catch (e) {
    status.textContent = "Parse error: " + e.message;
    return;
  }
  try {
    const res = await fetch("/api/solve", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ puzzle }),
    });
    const data = await res.json();
    if (data.error) {
      status.textContent = "Error: " + data.error;
      return;
    }
    renderInputPreview(document.getElementById("input-grid"), data.inputGrid);
    if (data.solved) {
      renderPipeGrid(
        document.getElementById("output-grid"),
        data.grid,
        data.inputGrid
      );
    } else {
      document.getElementById("output-grid").innerHTML = "";
    }
    status.textContent = data.solved ? "Solved" : "Unsolvable";
    const s = data.stats || {};
    document.getElementById("stats").textContent =
      `timeMs: ${data.timeMs}\nnodes: ${s.nodesExpanded ?? 0}\nbacktracks: ${s.backtracks ?? 0}\nmaxDepth: ${s.maxDepth ?? 0}`;
  } catch {
    status.textContent = "Request failed";
  }
}

function switchTab(name) {
  document.querySelectorAll(".tab").forEach((t) => {
    const on = t.dataset.tab === name;
    t.classList.toggle("active", on);
    t.setAttribute("aria-selected", on ? "true" : "false");
  });
  document.querySelectorAll(".panel").forEach((p) => {
    const on = p.id === `panel-${name}`;
    p.classList.toggle("active", on);
    p.hidden = !on;
  });
}

function init() {
  state.endpoints = emptyEndpoints();
  parsePuzzleText(DEFAULT_PUZZLE);
  document.getElementById("puzzle").value = DEFAULT_PUZZLE;
  document.getElementById("rows").value = state.rows;
  document.getElementById("cols").value = state.cols;

  renderPalette();
  renderEditor();

  document.querySelectorAll(".tab").forEach((tab) => {
    tab.addEventListener("click", () => {
      const name = tab.dataset.tab;
      if (name === "text") {
        document.getElementById("puzzle").value = toPuzzleText();
      }
      switchTab(name);
    });
  });

  document.querySelectorAll(".tool").forEach((btn) => {
    btn.addEventListener("click", () => {
      state.tool = btn.dataset.tool;
      document.querySelectorAll(".tool").forEach((t) => {
        t.classList.toggle("active", t === btn);
      });
    });
  });

  document.getElementById("resize-grid").addEventListener("click", () => {
    const rows = Math.min(12, Math.max(2, +document.getElementById("rows").value || 5));
    const cols = Math.min(12, Math.max(2, +document.getElementById("cols").value || 5));
    document.getElementById("rows").value = rows;
    document.getElementById("cols").value = cols;
    resetState(rows, cols, true);
    renderEditor();
    renderPalette();
  });

  document.getElementById("copy-text").addEventListener("click", () => {
    copyText(toPuzzleText(), document.getElementById("copy-status"));
  });

  document.getElementById("sync-to-text").addEventListener("click", () => {
    document.getElementById("puzzle").value = toPuzzleText();
    switchTab("text");
  });

  document.getElementById("load-text").addEventListener("click", () => {
    try {
      parsePuzzleText(document.getElementById("puzzle").value);
      document.getElementById("rows").value = state.rows;
      document.getElementById("cols").value = state.cols;
      renderEditor();
      renderPalette();
      switchTab("editor");
      document.getElementById("status").textContent = "Loaded into editor";
    } catch (e) {
      document.getElementById("status").textContent = "Parse error: " + e.message;
    }
  });

  document.getElementById("copy-textarea").addEventListener("click", () => {
    copyText(document.getElementById("puzzle").value, document.getElementById("status"));
  });

  document.getElementById("solve").addEventListener("click", solve);
}

init();
