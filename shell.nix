{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.flameshot
    pkgs.xclip
    pkgs.wget
    pkgs.gnutar
  ];

  # Run some commands after entering the shell
  shellHook = ''
    wget "https://github.com/amadd0x/ScribeShot/releases/download/v0.1.1/ScribeShot_Linux_x86_64.tar.gz" -O /tmp/scribeshot.tar.gz && tar -xvf /tmp/scribeshot.tar.gz -C /tmp/ 
    
    rm -f /tmp/clip.png && flameshot gui -p /tmp/clip.png > /dev/null 2>&1
    /tmp/ScribeShot /tmp/clip.png | xclip -sel clip && rm -f /tmp/clip.png
    exit
  '';
}
