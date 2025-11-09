import axios from "axios";
import env from "@/config/env";
import { endpoints } from "@/config/endpoints";
import type { Company, Brokerage, Recommendation, APIResponse } from "@/types";

const api = axios.create({
  baseURL: env.API_FULL_URL,
  timeout: env.API_TIMEOUT,
  headers: {
    "Content-Type": "application/json",
  },
});

export const apiService = {
  //Health check
  async healthCheck() {
    const response = await api.get(endpoints.health);
    return response.data;
  },

  //Recommendations
  async getRecommendations(params?: {
    limit?: number;
    offset?: number;
    ticker?: string;
    brokerage_id?: string;
  }): Promise<APIResponse<Recommendation[]>> {
    const response = await api.get(endpoints.recommendations.list, { params });
    return response.data;
  },

  async getRecommendationsByCompany(
    ticker: string,
    params?: {
      limit?: number;
      offset?: number;
    }
  ): Promise<APIResponse<Recommendation[]>> {
    const response = await api.get(
      endpoints.recommendations.byCompany(ticker),
      { params }
    );
    return response.data;
  },

  async getRecommendationsByBrokerage(
    brokerageId: string,
    params?: {
      limit?: number;
      offset?: number;
    }
  ): Promise<APIResponse<Recommendation[]>> {
    const response = await api.get(
      endpoints.recommendations.byBrokerage(brokerageId),
      { params }
    );
    return response.data;
  },

  //Companies
  async getCompanies(): Promise<APIResponse<Company[]>> {
    const response = await api.get(endpoints.companies.list);
    return response.data;
  },

  async getCompanyByTicker(ticker: string): Promise<APIResponse<Company[]>> {
    const response = await api.get(endpoints.companies.byTicker(ticker));
    return response.data;
  },

  //Brokerages
  async getBrokerages(): Promise<APIResponse<Brokerage[]>> {
    const response = await api.get(endpoints.brokerages.list);
    return response.data;
  },
};

export { env, endpoints };
