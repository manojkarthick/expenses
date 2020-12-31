with import <nixpkgs> {};

stdenv.mkDerivation rec {
  name = "env";
  env = buildEnv { name = name; paths = buildInputs; };
  buildInputs = [
	go_1_14
  ];

  shellHook = ''
    # export PS1="[nix@${name}:\W]\[$(tput sgr0)\] "
	export PS1="[\[\e[32m\]nix\[\e[m\]@\[\e[31m\]${name}\[\e[m\] \[\e[36m\]\W\[\e[m\]] "
    bind '"\e\e[C": forward-word'
    bind '"\e\e[D": backward-word'
  '';
}
