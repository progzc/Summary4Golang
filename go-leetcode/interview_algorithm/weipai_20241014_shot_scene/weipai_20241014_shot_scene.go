package weipai_20241014_shot_scene

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// 游乐场有个射击游戏，每次游戏前可花1块钱抽辅助道具。
// 抽到了不满意的道具可重复抽取，重复抽取时花费翻倍且只保留最新一个。
// 获得的道具会记在玩家身上，玩游戏时道具用掉。
// 玩家每天可有一次免费抽道具的机会，只有身上没有道具时抽才会免费。请编写程序帮摊主记录相关情况，
//
// 应有函数：查询下次抽道具的花费、抽道具、查询玩家当前拥有的道具、玩游戏。
// 数据记录用内部变量即可，一个玩家可能的行为序列
//【第一天->免费抽道具->花1块抽道具->玩游戏->花1块抽道具->第二天->花2块抽道具->玩游戏->免费抽道具->玩游戏 】

type ShotGame struct {
	Gifts    map[int]struct{}    // 道具
	DayRule  map[string]struct{} // 当天是否已抽
	LastCost int                 // 上次花费

	TotalCost       int            // 累计花费
	CostHistoryList []*CostHistory // 花费历史记录
	GameHistoryList []*GameHistory // 玩游戏历史记录
}

type CostHistory struct {
	Cost   int // 花费
	GiftID int // 道具 id
}

type GameHistory struct {
	GiftID   int       // 道具 id
	GameTime time.Time // 游戏消耗时间
}

//【第一天->免费抽道具->花1块抽道具->玩游戏->花1块抽道具->第二天->花2块抽道具->玩游戏->免费抽道具->玩游戏 】

// QryNextCost 查询下次抽道具的花费
func (s *ShotGame) QryNextCost() (cost int) {
	now := time.Now().Format("2006-01-02")
	if len(s.Gifts) == 0 {
		if _, ok := s.DayRule[now]; !ok { // 当天未抽取
			cost = 0
			return
		} else {
			cost = 1
			s.LastCost = 1
			return
		}
	}
	return s.LastCost * 2
}

// Lottery 抽道具
func (s *ShotGame) Lottery() (giftID int) {
	cost := s.QryNextCost()
	now := time.Now().Format("2006-01-02")
	giftID = int(rand.Int32())
	if cost == 0 {
		s.DayRule[now] = struct{}{}
	}
	s.Gifts[giftID] = struct{}{}
	s.TotalCost += cost
	s.LastCost = cost
	s.CostHistoryList = append(s.CostHistoryList, &CostHistory{Cost: cost, GiftID: giftID})
	return
}

// ExecGame 玩游戏
func (s *ShotGame) ExecGame(giftID int) (err error) {
	if _, ok := s.Gifts[giftID]; !ok {
		err = fmt.Errorf("giftID not found")
		return
	}
	s.GameHistoryList = append(s.GameHistoryList, &GameHistory{
		GiftID:   giftID,
		GameTime: time.Now(),
	})
	return
}

// QryAllGifts 查询所有的道具
func (s *ShotGame) QryAllGifts() (giftIDs []int) {
	for giftID := range s.Gifts {
		giftIDs = append(giftIDs, giftID)
	}
	return
}
