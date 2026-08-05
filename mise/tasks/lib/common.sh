# mise/tasks/lib/common.sh — shared helpers for Android/mobile tasks.
#
# Tasks source this with:  . "{{config_root}}/mise/tasks/lib/common.sh"
# (mise expands {{config_root}} before the script runs; each task is its own
# process, so helpers must be self-contained — no state between tasks).
#
# Keep everything bash/POSIX-ish and idempotent.

# irgo_detect_jdk17 [install_macos]
#   Detect a JDK 17 and export JAVA_HOME. Order:
#     1. existing JAVA_HOME (must have bin/java)
#     2. macOS /usr/libexec/java_home -v 17 (Temurin)
#     3. macOS brew openjdk@17 — installed when install_macos=1 and missing
#     4. Linux common JVM dirs (java-17-*, temurin-17*, msopenjdk-17*)
#     5. any `java` on PATH that actually runs (the macOS /usr/bin/java stub
#        fails `java -version`, so it is rejected)
#   Returns 0 if a usable JDK is available (JAVA_HOME exported, or java on
#   PATH), 1 otherwise. Never exits the caller.
irgo_detect_jdk17() {
  local install_macos="${1:-0}"
  if [ -z "$JAVA_HOME" ] || [ ! -x "$JAVA_HOME/bin/java" ]; then
    if [ -x /usr/libexec/java_home ] && JH="$(/usr/libexec/java_home -v 17 2>/dev/null)" && [ -n "$JH" ] && [ -x "$JH/bin/java" ]; then
      export JAVA_HOME="$JH"
    elif [ "$(uname -s)" = "Darwin" ]; then
      if [ "$install_macos" = "1" ] && ! brew list --versions openjdk@17 >/dev/null 2>&1; then
        echo "Installing openjdk@17 via brew (required by sdkmanager + gradle)..."
        brew install openjdk@17
      fi
      BP="$(brew --prefix openjdk@17 2>/dev/null || true)"
      if [ -n "$BP" ] && [ -d "$BP/libexec/openjdk.jdk/Contents/Home" ]; then
        export JAVA_HOME="$BP/libexec/openjdk.jdk/Contents/Home"
      fi
    elif [ "$(uname -s)" = "Linux" ]; then
      for d in /usr/lib/jvm/java-17-* /usr/lib/jvm/temurin-17* /usr/lib/jvm/msopenjdk-17*; do
        [ -x "$d/bin/java" ] && export JAVA_HOME="$d" && break
      done
    fi
  fi
  if [ -n "$JAVA_HOME" ] && [ -x "$JAVA_HOME/bin/java" ]; then
    export PATH="$JAVA_HOME/bin:$PATH"
    return 0
  fi
  if java -version >/dev/null 2>&1; then
    echo "Using java from PATH (JAVA_HOME not set)."
    return 0
  fi
  return 1
}

# irgo_locate_sdk_tools
#   Locate sdkmanager/avdmanager, preferring the SDK's own cmdline-tools so
#   the AVD is created against the right SDK root (brew-cask `command -v`
#   shims point at the cask's own root and produce broken AVDs with <build>
#   placeholders). Sets SDKM and AVDMGR (empty string when not found).
irgo_locate_sdk_tools() {
  SDKM="$(command -v sdkmanager 2>/dev/null || true)"
  AVDMGR="$(command -v avdmanager 2>/dev/null || true)"
  for p in "$ANDROID_HOME/cmdline-tools/latest/bin" \
           /opt/homebrew/share/android-commandlinetools/cmdline-tools/latest/bin \
           /usr/local/share/android-commandlinetools/cmdline-tools/latest/bin; do
    [ -z "$SDKM" ] && [ -x "$p/sdkmanager" ] && SDKM="$p/sdkmanager"
    [ -z "$AVDMGR" ] && [ -x "$p/avdmanager" ] && AVDMGR="$p/avdmanager"
  done
}
