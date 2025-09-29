export interface Company {
  id: string;
  ticker: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Brokerage {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Recommendation {
  id: string;
  company: Company;
  brokerage: Brokerage | null;
  target_from: number | null;
  target_to: number | null;
  rating_from: string;
  rating_to: string;
  action: string;
  time: string;
  created_at: string;
  updated_at: string;
}

export interface APIResponse<T> {
  success: boolean;
  data: T;
  error?: string;
  meta?: {
    total: number;
    limit: number;
    offset: number;
  };
}
