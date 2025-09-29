interface Environment {
    API_BASE_URL: string
    API_VERSION: string
    API_TIMEOUT: number
    API_FULL_URL: string
  }
  
  const env: Environment = {
    API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
    API_VERSION: import.meta.env.VITE_API_VERSION || 'v1',
    API_TIMEOUT: Number(import.meta.env.VITE_API_TIMEOUT) || 10000,
    get API_FULL_URL() {
      return `${this.API_BASE_URL}/api/${this.API_VERSION}`
    }
  }
  
  export default env