import { AppConfig } from './types';

// In this example, we get the config from environment variables.
// But you may define it directly

export default {
  appName: process.env['APP_NAME'],
  charonApiUrl: process.env['CHARON_API_URL'], // don't use trailing slash
  tinyMCEApiKey: 'lugruemgf9a7cb78atgaikhkaish8da7itsdaiusdausdhhy',
} as AppConfig;
