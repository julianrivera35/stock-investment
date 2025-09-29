import env from "./env";

export const endpoints ={
    base: env.API_BASE_URL,
    health: '/health',
    companies: {
        list: '/companies',
        byTicker: (ticker:string) => `/companies/${ticker}`
    },
    brokerages: {
        list: '/brokerages',
    },
    recommendations:{
        list: '/recommendations',
        byCompany: (ticker:string) =>`/recommendations/company/${ticker}`,
        byBrokerage: (brokerage:string) => `/recommendations/brokerage/${brokerage}`
    }
}

export const buildUrl = (endpoint:string) => `${endpoints.base}${endpoint}`