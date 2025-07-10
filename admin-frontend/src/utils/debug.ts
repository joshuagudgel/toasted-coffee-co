export const debugFetch = async (url: string, options: RequestInit) => {
  console.log(`Fetching ${url}`, options);
  try {
    const response = await fetch(url, options);
    console.log(`Response from ${url}:`, response.status);
    return response;
  } catch (error) {
    console.error(`Error fetching ${url}:`, error);
    throw error;
  }
};