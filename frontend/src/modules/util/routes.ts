export function generateUrlWithParams(url: string, params: { [key: string]: any }) {
  let urlWithParams = url;
  for (const key in params) {
    urlWithParams = urlWithParams.replace(`:${key}`, params[key]);
  }
  return urlWithParams;
}