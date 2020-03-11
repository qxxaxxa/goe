package main

import (
	"fmt"
	"regexp"
)

func main() {
	var chLoGrpReg, _ = regexp.Compile(`\[[^\[]+[汉漢]+[^]]*]`)
	s := "(BanG Dreamer's Party! 8th STAGE) [Tohonifun (Chado)] Hana no Iro wa | 花的顏色 (BanG Dream!) [Chinese] [EZR個人漢化][Kihanabatake (Kihana)] Ima, Soko ni Iru Kimi e. | 致如今、身在此處的妳 (Touhou Project) [Chinese] [魔戀漢化組][Macchadokoro (Warashibe)] Anal Megumi wa Sukidarake (Amano Megumi ha Sukidarake!)] | 屁眼惠破綻百出 [Chinese] [天帝哥個人漢化][Kyockcho] Shiori Panic (COMIC BAVEL 2019-06) [Chinese] [無邪気漢化組] [Digital](C96) [OrangeMaru (YD)] Chaldea Maid #Mash (Fate/Grand Order) [Chinese] [空気系☆漢化] [Decensored][Ne. (Shiromitsu Daiya)] Anju to Mazareba Amai Mitsu [Chinese] [瑞树汉化组] [Digital][Horie Tankei] Haha no Himitsu | Secret of Mother Ch. 1-4 [Chinese] [官能战士个人汉化](C97) [Yuuki Nyuugyou (Yuuki Shin)] Makoto no Ai (Princess Connect! Re:Dive) [Chinese] [廢欲加速漢化](BanG Dreamer's Party! 2nd STAGE) [majicalcarca (yae)] Hajimete no o-tomari. (BanG Dream!) [Chinese] [猫在汉化](C93) [Chilly polka (Suimya)] ChillypoRoom WMvol.11 [Chinese] [绅士仓库汉化][Same Manma] Juuden max!(COMIC Kairakuten 2018-8) [2020年1月21日漢化][IRODORI (SOYOSOYO)] Natsu no Yuuutsutsu [Chinese] [lolipoi汉化组] [Digital](C87) [Kanimiso-tei (Kusatsu Terunyo)] Azusa-San Maji Tekireiki (The Idolm@ster) [Chinese] [个人不完全渣渣汉化]"
	allString := chLoGrpReg.FindAllString(s, -1)
	for _, s2 := range allString {
		fmt.Println(s2)
	}
}
