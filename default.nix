{ lib, buildGoModule, i3status }:

buildGoModule {
  pname = "i3statuspp";
  version = "0.0.1";

  src = lib.cleanSource ./.;

  vendorSha256 = "0g07spr483l0cpw73dls91nsz5mydgijdcff65fmmac3cqhp3p9i";

  buildFlagsArray = [
    "-ldflags=-s -w"
    "-extldflags \"-static\""
  ];

  meta = with lib; {
    description = "Extend i3status functionality";
    homepage = "https://github.com/michaeladler/i3statuspp";
    license = licenses.asl20;
  };
}
