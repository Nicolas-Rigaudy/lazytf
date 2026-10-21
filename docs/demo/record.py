#!/usr/bin/env python3
"""Record the README screenshots: run lazytf on docs/demo/sample, press keys, dump screens.

  python3 docs/demo/record.py LAZYTF_BINARY OUT_DIR [--cols 140] [--rows 30] [--plain-symbols]

Copies the sample repo to OUT_DIR/work (so state and .terraform never land in
docs/), runs the binary in a pseudo-terminal with HOME pointed at OUT_DIR/home
(so your own config and caches are not used), feeds it the key presses below and
writes the screen after each step as OUT_DIR/<name>.ans, ANSI text that any
ANSI-to-PNG renderer turns into an image. Needs Python 3, pyte and terraform.
"""
import argparse
import fcntl
import os
import pty
import select
import shutil
import struct
import sys
import termios
import time

import pyte

HERE = os.path.dirname(os.path.abspath(__file__))

# (keys to send, text that must be on screen before the shot, shot name or None)
STEPS = [
    ("", "app", "1-discover"),
    ("\r", "dev", "2-environments"),
    ("\r", "not initialized", None),
    ("y", "Confirm Terraform Init", "3-init-confirm"),
    ("y", "successfully initialized", "4-init"),
    ("p", "Plan Summary", None),
    ("\t", "Plan Summary", None),  # focus the output panel
    ("\n", "module.api", None),  # ctrl+j: next row
    ("\r", "terraform_data.service", "5-plan-summary"),
    ("a", "Use saved plan", "6-apply-saved-plan"),
    ("\r", "Apply complete", "7-apply"),
]

# Colour emoji drawn as one-cell text symbols, for --plain-symbols.
PLAIN = {"\u2705": "\u2714", "\u274c": "\u2718", "\U0001f4cb": "\u2261"}

NAMED = {
    "black": 0, "red": 1, "green": 2, "brown": 3, "blue": 4, "magenta": 5, "cyan": 6, "white": 7,
    "brightblack": 8, "brightred": 9, "brightgreen": 10, "brightbrown": 11,
    "brightblue": 12, "brightmagenta": 13, "brightcyan": 14, "brightwhite": 15,
}


def sgr_color(color, fg):
    if color == "default":
        return "39" if fg else "49"
    if color in NAMED:
        n = NAMED[color]
        base = (30 if fg else 40) if n < 8 else (90 if fg else 100)
        return str(base + n % 8)
    r, g, b = int(color[0:2], 16), int(color[2:4], 16), int(color[4:6], 16)
    return f"{38 if fg else 48};2;{r};{g};{b}"


def screen_to_ansi(screen, plain=False, masks=()):
    """One character per cell: the right half of a wide character becomes a space.

    Each (old, new) in masks replaces old with new and pads the rest of that word
    with spaces to the same width, so the shots show ~/infra instead of the
    recording machine's path.
    """
    lines = []
    for y in range(screen.lines):
        row, out, last = screen.buffer[y], [], None
        chars = [row[x].data or "\0" for x in range(screen.columns)]
        for old, new in masks:
            text = "".join(chars)
            i = text.find(old)
            while i >= 0:
                j = i + len(old)  # end of the word old starts, e.g. /app after the path
                while j < len(chars) and chars[j] != " ":
                    j += 1
                chars[i:j] = list(new) + chars[i + len(old):j] + [" "] * (len(old) - len(new))
                text = "".join(chars)
                i = text.find(old, i + len(new))
        for x in range(screen.columns):
            ch = row[x]
            attrs = (ch.fg, ch.bg, ch.bold, ch.italics, ch.underscore, ch.reverse)
            if attrs != last:
                codes = ["0", sgr_color(ch.fg, True), sgr_color(ch.bg, False)]
                codes += [c for c, on in zip("1347", attrs[2:]) if on]
                out.append("\x1b[" + ";".join(codes) + "m")
                last = attrs
            data = chars[x] if chars[x] != "\0" else " "
            out.append(PLAIN.get(data, data) if plain else data)
        lines.append("".join(out) + "\x1b[0m")
    return "\n".join(lines) + "\n"


def screen_text(screen):
    return "\n".join(screen.display)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("binary")
    ap.add_argument("out")
    ap.add_argument("--cols", type=int, default=140)
    ap.add_argument("--rows", type=int, default=30)
    ap.add_argument("--timeout", type=float, default=60)
    ap.add_argument("--plain-symbols", action="store_true",
                    help="replace colour emoji with text symbols, for renderers without an emoji font")
    a = ap.parse_args()

    out = os.path.abspath(a.out)
    work, home = os.path.join(out, "work"), os.path.join(out, "home")
    shutil.rmtree(work, ignore_errors=True)
    shutil.rmtree(home, ignore_errors=True)
    shutil.copytree(os.path.join(HERE, "sample"), work)
    os.makedirs(home)

    env = dict(os.environ, HOME=home, TERM="xterm-256color", COLORTERM="truecolor")
    env.pop("XDG_CONFIG_HOME", None)
    env.pop("XDG_CACHE_HOME", None)
    binary = os.path.abspath(a.binary)

    pid, fd = pty.fork()
    if pid == 0:
        os.chdir(work)
        os.execve(binary, [binary], env)
    fcntl.ioctl(fd, termios.TIOCSWINSZ, struct.pack("HHHH", a.rows, a.cols, 0, 0))

    screen = pyte.Screen(a.cols, a.rows)
    stream = pyte.ByteStream(screen)

    def pump(seconds):
        end = time.time() + seconds
        while True:
            left = end - time.time()
            if left <= 0:
                return
            r, _, _ = select.select([fd], [], [], left)
            if r:
                try:
                    data = os.read(fd, 65536)
                except OSError:
                    return
                stream.feed(data)

    try:
        pump(1.5)
        for keys, expect, name in STEPS:
            if keys:
                os.write(fd, keys.encode())
            deadline = time.time() + a.timeout
            while expect not in screen_text(screen):
                if time.time() > deadline:
                    sys.exit(f"timed out waiting for {expect!r} after {keys!r}:\n{screen_text(screen)}")
                pump(0.2)
            pump(1.0)  # let the frame settle
            if name:
                with open(os.path.join(out, name + ".ans"), "w") as f:
                    f.write(screen_to_ansi(screen, a.plain_symbols, [(work, "~/infra")]))
                print("wrote", name)
    finally:
        os.write(fd, b"\x03")
        pump(0.5)
        os.close(fd)
        os.waitpid(pid, 0)


if __name__ == "__main__":
    main()
