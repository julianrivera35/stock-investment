import type { Recommendation } from "@/types";

export interface StockScore {
  ticker: string;
  companyName: string;
  totalScore: number;
  metrics: {
    recommendationCount: number;
    upgradeRatio: number;
    averageTargetUpside: number;
    analystConsensus: number;
    recencyScore: number;
    brokerageQuality: number;
  };
  details: {
    totalRecommendations: number;
    upgrades: number;
    downgrades: number;
    reiterations: number;
    averageTargetPrice: number;
    currentTargetPrice: number;
    uniqueBrokerages: number;
    daysFromLastRecommendation: number;
  };
  recommendation: "Strong Buy" | "Buy" | "Hold" | "Sell" | "Strong Sell";
}

function getPremiumBrokerages(recommendations: Recommendation[]): string[] {
  const brokerageStats: Record<
    string,
    {
      name: string;
      totalRecommendations: number;
      accuracyScore: number;
      uniqueCompanies: number;
    }
  > = {};

  recommendations.forEach((rec) => {
    const name = rec.brokerage?.name;
    if (!name) return;

    if (!brokerageStats[name]) {
      brokerageStats[name] = {
        name,
        totalRecommendations: 0,
        accuracyScore: 0,
        uniqueCompanies: 0,
      };
    }

    brokerageStats[name].totalRecommendations++;
  });

  // Calculate unique companies per brokerage
  const brokerageCompanies: Record<string, Set<string>> = {};
  recommendations.forEach((rec) => {
    const brokerage = rec.brokerage?.name;
    const company = rec.company.ticker;
    if (!brokerage) return;

    if (!brokerageCompanies[brokerage]) {
      brokerageCompanies[brokerage] = new Set();
    }
    brokerageCompanies[brokerage].add(company);
  });

  // Calculate accuracy scores (upgrade ratio per brokerage)
  Object.keys(brokerageStats).forEach((brokerageName) => {
    const brokerageRecs = recommendations.filter(
      (r) => r.brokerage?.name === brokerageName
    );
    const upgrades = brokerageRecs.filter(
      (r) =>
        r.action.toLowerCase().includes("raised") ||
        r.action.toLowerCase().includes("upgrade")
    ).length;

    brokerageStats[brokerageName].accuracyScore =
      brokerageRecs.length > 0 ? (upgrades / brokerageRecs.length) * 100 : 0;

    brokerageStats[brokerageName].uniqueCompanies =
      brokerageCompanies[brokerageName]?.size || 0;
  });

  // Determine premium brokerages based on data
  // Premium = high volume + good accuracy + diverse coverage
  const sortedBrokerages = Object.values(brokerageStats)
    .filter((b) => b.totalRecommendations >= 5)
    .map((brokerage) => ({
      ...brokerage,
      qualityScore:
        Math.min(brokerage.totalRecommendations / 20, 1) * 40 +
        brokerage.accuracyScore * 0.4 +
        Math.min(brokerage.uniqueCompanies / 10, 1) * 20,
    }))
    .sort((a, b) => b.qualityScore - a.qualityScore);

  const premiumCount = Math.max(3, Math.floor(sortedBrokerages.length * 0.3));
  return sortedBrokerages.slice(0, premiumCount).map((b) => b.name);
}

export function analyzeStockInvestments(
  recommendations: Recommendation[]
): StockScore[] {
  // Group recommendations by company
  const stockGroups = groupRecommendationsByStock(recommendations);

  // Get premium brokerages from actual data
  const premiumBrokerages = getPremiumBrokerages(recommendations);
  console.log("Data-driven premium brokerages:", premiumBrokerages);

  // Analyze each stock
  const stockScores = Object.entries(stockGroups).map(([ticker, recs]) =>
    analyzeStock(ticker, recs, premiumBrokerages)
  );

  // Sort by total score (highest first)
  return stockScores.sort((a, b) => b.totalScore - a.totalScore);
}

function groupRecommendationsByStock(
  recommendations: Recommendation[]
): Record<string, Recommendation[]> {
  return recommendations.reduce((groups, rec) => {
    const ticker = rec.company.ticker;
    if (!groups[ticker]) {
      groups[ticker] = [];
    }
    groups[ticker].push(rec);
    return groups;
  }, {} as Record<string, Recommendation[]>);
}

function analyzeStock(
  ticker: string,
  recommendations: Recommendation[],
  premiumBrokerages: string[]
): StockScore {
  const companyName = recommendations[0]?.company.name || ticker;

  // Calculate basic metrics
  const upgrades = recommendations.filter(
    (r) =>
      r.action.toLowerCase().includes("raised") ||
      r.action.toLowerCase().includes("upgrade")
  );
  const downgrades = recommendations.filter(
    (r) =>
      r.action.toLowerCase().includes("lowered") ||
      r.action.toLowerCase().includes("downgrade")
  );
  const reiterations = recommendations.filter((r) =>
    r.action.toLowerCase().includes("reiterat")
  );

  // Get valid target prices
  const validTargets = recommendations
    .filter((r) => r.target_to && r.target_to > 0)
    .map((r) => r.target_to!);

  const averageTargetPrice =
    validTargets.length > 0
      ? validTargets.reduce((sum, price) => sum + price, 0) /
        validTargets.length
      : 0;

  const currentTargetPrice =
    validTargets.length > 0 ? Math.max(...validTargets) : 0;

  // Calculate target upside
  const minTargetPrice =
    validTargets.length > 0 ? Math.min(...validTargets) : 0;
  const averageTargetUpside =
    minTargetPrice > 0
      ? ((averageTargetPrice - minTargetPrice) / minTargetPrice) * 100
      : 0;

  // Get unique brokerages
  const uniqueBrokerages = [
    ...new Set(recommendations.map((r) => r.brokerage?.name).filter(Boolean)),
  ].length;

  // Calculate days from last recommendation
  const latestRecommendation = recommendations
    .map((r) => new Date(r.time))
    .sort((a, b) => b.getTime() - a.getTime())[0];

  const daysFromLastRecommendation = latestRecommendation
    ? Math.floor(
        (Date.now() - latestRecommendation.getTime()) / (1000 * 60 * 60 * 24)
      )
    : 999;

  // Calculate metrics (0-100 scale)
  const metrics = {
    recommendationCount: Math.min((recommendations.length / 20) * 100, 100),
    upgradeRatio:
      recommendations.length > 0
        ? (upgrades.length / recommendations.length) * 100
        : 0,
    averageTargetUpside: Math.min(Math.max(averageTargetUpside, 0), 100),
    analystConsensus: calculateConsensusScore(
      upgrades.length,
      downgrades.length,
      reiterations.length
    ),
    recencyScore: Math.max(100 - daysFromLastRecommendation * 2, 0),
    brokerageQuality: calculateBrokerageQualityScore(
      recommendations,
      premiumBrokerages
    ),
  };

  // Weight the scores
  const weights = {
    recommendationCount: 0.15,
    upgradeRatio: 0.25,
    averageTargetUpside: 0.2,
    analystConsensus: 0.2,
    recencyScore: 0.1,
    brokerageQuality: 0.1,
  };

  const totalScore = Object.entries(metrics).reduce((total, [key, value]) => {
    return total + value * weights[key as keyof typeof weights];
  }, 0);

  const details = {
    totalRecommendations: recommendations.length,
    upgrades: upgrades.length,
    downgrades: downgrades.length,
    reiterations: reiterations.length,
    averageTargetPrice,
    currentTargetPrice,
    uniqueBrokerages,
    daysFromLastRecommendation,
  };

  return {
    ticker,
    companyName,
    totalScore: Math.round(totalScore * 100) / 100,
    metrics,
    details,
    recommendation: getRecommendationLabel(totalScore),
  };
}

function calculateConsensusScore(
  upgrades: number,
  downgrades: number,
  reiterations: number
): number {
  const total = upgrades + downgrades + reiterations;
  if (total === 0) return 0;

  // Strong consensus
  const maxCategory = Math.max(upgrades, downgrades, reiterations);
  const consensusRatio = maxCategory / total;

  // Bonus for upgrade consensus
  const isUpgradeConsensus = upgrades === maxCategory;
  const bonus = isUpgradeConsensus ? 20 : 0;

  return Math.min(consensusRatio * 80 + bonus, 100);
}

function calculateBrokerageQualityScore(
  recommendations: Recommendation[],
  premiumBrokerages: string[]
): number {
  const brokerageNames = recommendations
    .map((r) => r.brokerage?.name)
    .filter(Boolean);

  if (brokerageNames.length === 0) return 0;

  const premiumCount = brokerageNames.filter((name) =>
    premiumBrokerages.some((premium) => name!.includes(premium))
  ).length;

  const premiumRatio = premiumCount / brokerageNames.length;
  const diversityBonus = Math.min(brokerageNames.length / 5, 1) * 20;

  return Math.min(premiumRatio * 80 + diversityBonus, 100);
}

function getRecommendationLabel(score: number): StockScore["recommendation"] {
  if (score >= 80) return "Strong Buy";
  if (score >= 65) return "Buy";
  if (score >= 40) return "Hold";
  if (score >= 20) return "Sell";
  return "Strong Sell";
}

export function getTopInvestmentPicks(
  recommendations: Recommendation[],
  count: number = 5
): StockScore[] {
  return analyzeStockInvestments(recommendations).slice(0, count);
}
