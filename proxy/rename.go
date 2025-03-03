package proxies

import (
	"regexp"
	"strconv"
	"sync/atomic"
)

// Counter 用于存储各个地区的计数
type Counter struct {
	// 香港
	hk int32
	// 台湾
	tw int32
	// 美国
	us int32
	// 新加坡
	sg int32
	// 日本
	jp int32
	// 英国
	uk int32
	// 加拿大
	ca int32
	// 澳大利亚
	au int32
	// 德国
	de int32
	// 法国
	fr int32
	// 荷兰
	nl int32
	// 俄罗斯
	ru int32
	// 匈牙利
	hu int32
	// 乌克兰
	ua int32
	// 波兰
	pl int32
	// 韩国
	kr int32
	// 亚太地区
	ap int32
	// 伊朗
	ir int32
	// 意大利
	it int32
	// 芬兰
	fi int32
	// 柬埔寨
	kh int32
	// 巴西
	br int32
	// 印度
	in int32
	// 阿拉伯酋长国
	ae int32
	// 瑞士
	ch int32
	// 其他
	other int32
}

var counter Counter

// Reset 重置所有计数器为0
func ResetRenameCounter() {
	counter = Counter{}
}

func Rename(name string) string {
	// 香港
	if regexp.MustCompile(`(?i)(hk|港|hongkong|hong kong)`).MatchString(name) {
		atomic.AddInt32(&counter.hk, 1)
		return "🇭🇰香港" + strconv.Itoa(int(atomic.LoadInt32(&counter.hk)))
	}
	// 台湾
	if regexp.MustCompile(`(?i)(tw|台|taiwan|tai wen)`).MatchString(name) {
		atomic.AddInt32(&counter.tw, 1)
		return "🇹🇼台湾" + strconv.Itoa(int(atomic.LoadInt32(&counter.tw)))
	}
	// 美国
	if regexp.MustCompile(`(?i)(us|美|united states|america)`).MatchString(name) {
		atomic.AddInt32(&counter.us, 1)
		return "🇺🇸美国" + strconv.Itoa(int(atomic.LoadInt32(&counter.us)))
	}
	// 新加坡
	if regexp.MustCompile(`(?i)(sg|新|singapore|狮城)`).MatchString(name) {
		atomic.AddInt32(&counter.sg, 1)
		return "🇸🇬新加坡" + strconv.Itoa(int(atomic.LoadInt32(&counter.sg)))
	}
	// 日本
	if regexp.MustCompile(`(?i)(jp|日|japan)`).MatchString(name) {
		atomic.AddInt32(&counter.jp, 1)
		return "🇯🇵日本" + strconv.Itoa(int(atomic.LoadInt32(&counter.jp)))
	}
	// 英国
	if regexp.MustCompile(`(?i)(uk|英|united kingdom|britain|gb)`).MatchString(name) {
		atomic.AddInt32(&counter.uk, 1)
		return "🇬🇧英国" + strconv.Itoa(int(atomic.LoadInt32(&counter.uk)))
	}
	// 加拿大
	if regexp.MustCompile(`(?i)(ca|加|canada)`).MatchString(name) {
		atomic.AddInt32(&counter.ca, 1)
		return "🇨🇦加拿大" + strconv.Itoa(int(atomic.LoadInt32(&counter.ca)))
	}
	// 澳大利亚
	if regexp.MustCompile(`(?i)(au|澳|australia)`).MatchString(name) {
		atomic.AddInt32(&counter.au, 1)
		return "🇦🇺澳大利亚" + strconv.Itoa(int(atomic.LoadInt32(&counter.au)))
	}
	// 德国
	if regexp.MustCompile(`(?i)(de|德|germany|deutschland)`).MatchString(name) {
		atomic.AddInt32(&counter.de, 1)
		return "🇩🇪德国" + strconv.Itoa(int(atomic.LoadInt32(&counter.de)))
	}
	// 法国
	if regexp.MustCompile(`(?i)(fr|法|france)`).MatchString(name) {
		atomic.AddInt32(&counter.fr, 1)
		return "🇫🇷法国" + strconv.Itoa(int(atomic.LoadInt32(&counter.fr)))
	}
	// 荷兰
	if regexp.MustCompile(`(?i)(nl|荷|netherlands)`).MatchString(name) {
		atomic.AddInt32(&counter.nl, 1)
		return "🇳🇱荷兰" + strconv.Itoa(int(atomic.LoadInt32(&counter.nl)))
	}
	// 俄罗斯
	if regexp.MustCompile(`(?i)(ru|俄|russia)`).MatchString(name) {
		atomic.AddInt32(&counter.ru, 1)
		return "🇷🇺俄罗斯" + strconv.Itoa(int(atomic.LoadInt32(&counter.ru)))
	}
	// 匈牙利
	if regexp.MustCompile(`(?i)(hu|匈|hungary)`).MatchString(name) {
		atomic.AddInt32(&counter.hu, 1)
		return "🇭🇺匈牙利" + strconv.Itoa(int(atomic.LoadInt32(&counter.hu)))
	}
	// 乌克兰
	if regexp.MustCompile(`(?i)(ua|乌|ukraine)`).MatchString(name) {
		atomic.AddInt32(&counter.ua, 1)
		return "🇺🇦乌克兰" + strconv.Itoa(int(atomic.LoadInt32(&counter.ua)))
	}
	// 波兰
	if regexp.MustCompile(`(?i)(pl|波|poland)`).MatchString(name) {
		atomic.AddInt32(&counter.pl, 1)
		return "🇵🇱波兰" + strconv.Itoa(int(atomic.LoadInt32(&counter.pl)))
	}
	// 韩国
	if regexp.MustCompile(`(?i)(kr|韩|korea)`).MatchString(name) {
		atomic.AddInt32(&counter.kr, 1)
		return "🇰🇷韩国" + strconv.Itoa(int(atomic.LoadInt32(&counter.kr)))
	}
	// 亚太地区
	if regexp.MustCompile(`(?i)(ap|亚太|asia)`).MatchString(name) {
		atomic.AddInt32(&counter.ap, 1)
		return "🌏亚太地区" + strconv.Itoa(int(atomic.LoadInt32(&counter.ap)))
	}
	// 伊朗
	if regexp.MustCompile(`(?i)(ir|伊|iran)`).MatchString(name) {
		atomic.AddInt32(&counter.ir, 1)
		return "🇮🇷伊朗" + strconv.Itoa(int(atomic.LoadInt32(&counter.ir)))
	}
	// 意大利
	if regexp.MustCompile(`(?i)(it|意|italy)`).MatchString(name) {
		atomic.AddInt32(&counter.it, 1)
		return "🇮🇹意大利" + strconv.Itoa(int(atomic.LoadInt32(&counter.it)))
	}
	// 芬兰
	if regexp.MustCompile(`(?i)(fi|芬|finland)`).MatchString(name) {
		atomic.AddInt32(&counter.fi, 1)
		return "🇫🇮芬兰" + strconv.Itoa(int(atomic.LoadInt32(&counter.fi)))
	}
	// 柬埔寨
	if regexp.MustCompile(`(?i)(kh|柬|cambodia)`).MatchString(name) {
		atomic.AddInt32(&counter.kh, 1)
		return "🇰🇭柬埔寨" + strconv.Itoa(int(atomic.LoadInt32(&counter.kh)))
	}
	// 巴西
	if regexp.MustCompile(`(?i)(br|巴|brazil)`).MatchString(name) {
		atomic.AddInt32(&counter.br, 1)
		return "🇧🇷巴西" + strconv.Itoa(int(atomic.LoadInt32(&counter.br)))
	}
	// 印度
	if regexp.MustCompile(`(?i)(in|印|india)`).MatchString(name) {
		atomic.AddInt32(&counter.in, 1)
		return "🇮🇳印度" + strconv.Itoa(int(atomic.LoadInt32(&counter.in)))
	}
	// 阿拉伯酋长国
	if regexp.MustCompile(`(?i)(ae|阿|uae|阿拉伯酋长国)`).MatchString(name) {
		atomic.AddInt32(&counter.ae, 1)
		return "🇦🇪阿拉伯酋长国" + strconv.Itoa(int(atomic.LoadInt32(&counter.ae)))
	}
	// 瑞士
	if regexp.MustCompile(`(?i)(ch|瑞|switzerland)`).MatchString(name) {
		atomic.AddInt32(&counter.ch, 1)
		return "🇨🇭瑞士" + strconv.Itoa(int(atomic.LoadInt32(&counter.ch)))
	}
	// 其他
	atomic.AddInt32(&counter.other, 1)
	return "🌀其他" + strconv.Itoa(int(atomic.LoadInt32(&counter.other))) + "-" + name
}
