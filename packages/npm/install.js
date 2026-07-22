const { existsSync, mkdirSync, chmodSync, createWriteStream } = require('fs');
const { join, homedir } = require('path');
const { platform, arch } = require('os');
const https = require('https');
const { createGunzip } = require('zlib');
const { spawnSync } = require('child_process');

const CACHE_DIR = join(homedir(), '.game-dev-cli', 'bin');

function platformTag() {
  const m = { darwin: 'darwin', linux: 'linux', win32: 'windows' };
  return m[platform()] || platform();
}

function archTag() {
  const m = { x64: 'amd64', arm64: 'arm64' };
  return m[arch()] || arch();
}

function binaryName() {
  const ext = platform() === 'win32' ? '.exe' : '';
  return `game-dev-cli${ext}`;
}

function tarballName() {
  return `game-dev-cli-${platformTag()}-${archTag()}.tar.gz`;
}

// ponytail: replace with go-install or a dedicated downloader binary for cross-platform safety.
function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = createWriteStream(dest);
    https.get(url, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        file.close();
        download(res.headers.location, dest).then(resolve).catch(reject);
        return;
      }
      if (res.statusCode !== 200) {
        file.close();
        reject(new Error(`HTTP ${res.statusCode} — ${url}`));
        return;
      }
      res.pipe(file);
      file.on('finish', () => { file.close(); resolve(); });
    }).on('error', reject);
  });
}

function ensureBinary() {
  const bin = binaryName();
  const path = join(CACHE_DIR, bin);

  if (existsSync(path)) {
    // lazy: no version check, trusting npm/npx's cache invalidation on reinstall
    return path;
  }

  mkdirSync(CACHE_DIR, { recursive: true });

  const tag = process.env.GAME_DEV_CLI_VERSION || 'latest';
  const url = tag === 'latest'
    ? `https://github.com/8liang/game-dev-cli/releases/latest/download/${tarballName()}`
    : `https://github.com/8liang/game-dev-cli/releases/download/${tag}/${tarballName()}`;

  const tmpTar = join(CACHE_DIR, tarballName());
  download(url, tmpTar);

  // extract tar.gz: game-dev-cli<platform>-<arch> → game-dev-cli
  const { status, stderr } = spawnSync('tar', ['xzf', tmpTar, '--strip-components=0'], {
    cwd: CACHE_DIR,
    stdio: ['ignore', 'inherit', 'pipe'],
  });
  if (status !== 0) {
    throw new Error(`tar extract failed: ${stderr.toString().trim()}`);
  }

  // the tarball contains the full name, rename to plain binary name
  const exactName = `game-dev-cli-${platformTag()}-${archTag()}`;
  const exactPath = join(CACHE_DIR, exactName);
  if (existsSync(exactPath)) {
    spawnSync('mv', [exactPath, path]);
  }

  chmodSync(path, 0o755);
  try { spawnSync('rm', ['-f', tmpTar]); } catch {}

  return path;
}

function install() {
  try {
    return ensureBinary();
  } catch (err) {
    console.error('game-dev-cli 安装失败:', err.message);
    process.exit(1);
  }
}

module.exports = { install };

// allow standalone run: node install.js
if (require.main === module) {
  install();
}
