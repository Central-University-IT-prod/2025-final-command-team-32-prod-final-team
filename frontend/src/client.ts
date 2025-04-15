import { toast } from "sonner";

export type JWT = {
  id: string;
  username: string;
};

export type Cinema = {
  id: string;
  name: string;
  name_original?: string;
  year?: number;
  age_rating?: number;
  duration_minutes?: number;
  poster_url?: string;
  description?: string;
  genres?: string[];
  actors?: string[];
  rating?: number;
  user_rate?: number;
  user_like_status?: number;
};

export type Couch = {
  id: string;
  name: string;
  users: string[];
  author: string;
};

const API_BASE = "https://prod-team-32-n26k57br.REDACTED/api/v1";

async function request<T>(
  endpoint: string,
  method: string = "GET",
  body?: any,
  token?: string,
): Promise<T> {
  const headers: HeadersInit = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${API_BASE}${endpoint}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!res.ok) {
    //TODO: remove errors before deadline
    const errorText = await res.text();
    toast.error(`HTTP Error ${res.status}: ${errorText}`);
    throw new Error(`HTTP Error: ${res.status} - ${errorText}`);
  }

  try {
    return (await res.json()) as T;
  } catch {
    return {} as T;
  }
}

export type User = {
  id: string;
  username: string;
};
export const deleteFilm = (id: string, token: string) =>
  request<void>(`/admin/films/${id}`, "DELETE", undefined, token);

export const updateFilm = (
  id: string,
  data: {
    name: string;
    description?: string;
    poster_url?: string;
    genres?: string[];
  },
  token: string,
) => request<Cinema>(`/admin/films/${id}`, "PUT", data, token);

export const searchUsers = (query: string) =>
  request<User[]>(`/users/search?query=${query}`);

// Users API
export const signIn = (username: string, password: string) =>
  request<{ token: JWT }>("/users/sign-in", "POST", { username, password });

export const signUp = (username: string, password: string) =>
  request<{ token: JWT }>("/users/sign-up", "POST", { username, password });

export const uploadPoster = (id: string, file: File, token: string) => {
  const formData = new FormData();
  formData.append("uploadfile", file);
  return request<{ url: string }>(
    `/films/${id}/picture`,
    "POST",
    formData,
    token,
  );
};

// Cinema API
export const getGenres = () => request<string[]>(`/films/genres`);

export const getFilm = (id: string) => request<Cinema>(`/films/${id}`);

export const getFilmFeed = (limit: number = 10, token: string) =>
  request<Cinema[]>(`/films/feed?limit=${limit}`, "GET", null, token);

export const getCouchFeed = (id: string, limit: number = 10) =>
  request<Cinema[]>(`/couches/${id}/feed?limit=${limit}`, "GET");

export const seenCouch = (ids: string[], couch_id: string) =>
  request<void>(`/couches/${couch_id}/views/bulk`, "POST", ids);

export const searchFilm = (query: string, tags?: string[]) =>
  request<Cinema[]>(
    `/films/search?query=${query}${tags && "&tags=" + tags.join(",")}`,
  );

export const getPopularFilms = (limit: number = 10, offset: number = 0) =>
  request<Cinema[]>(`/films/popular?limit=${limit}&offset=${offset}`);

export const seen = (ids: string[], token: string) =>
  request<void>(`/films/views/bulk`, "POST", ids, token);

export const addFilm = (
  name: string,
  description?: string,
  poster_url?: string,
  genres?: string[],
  token?: string,
) =>
  request<{ id: string }>(
    "/films",
    "POST",
    { name, poster_url, genres, description },
    token,
  );

export const addAdminFilm = (
  name: string,
  description?: string,
  poster_url?: string,
  genres?: string[],
  token?: string,
) =>
  request<{ id: string }>(
    "/admin/films/",
    "POST",
    { name, poster_url, genres, description },
    token,
  );

// Likes API
export const likeFilm = (id: string, token: string) =>
  request<void>(`/films/${id}/like`, "POST", undefined, token);

export const dislikeFilm = (id: string, token: string) =>
  request<void>(`/films/${id}/dislike`, "POST", undefined, token);

export const removeLike = (id: string, token: string) =>
  request<void>(`/films/${id}/like`, "DELETE", undefined, token);

export const removeDislike = (id: string, token: string) =>
  request<void>(`/films/${id}/dislike`, "DELETE", undefined, token);

export const likeFilmCouch = (id: string, group: string, token: string) =>
  request<void>(`/couches/${group}/films/${id}/like`, "POST", undefined, token);

export const dislikeFilmCouch = (id: string, group: string, token: string) =>
  request<void>(
    `/couches/${group}/films/${id}/dislike`,
    "POST",
    undefined,
    token,
  );

export const removeLikeCouch = (id: string, group: string, token: string) =>
  request<void>(
    `/couches/${group}/films/${id}/like`,
    "DELETE",
    undefined,
    token,
  );

export const removeDislikeCouch = (id: string, group: string, token: string) =>
  request<void>(
    `/couches/${group}/films/${id}/dislike`,
    "DELETE",
    undefined,
    token,
  );

export const rateFilm = (id: string, rate: number, token: string) =>
  request<void>(`/films/${id}/rate`, "POST", { rate }, token);

// Plans API
export const getPlans = (token: string, limit: number = 30) =>
  request<Cinema[]>(`/plans?limit=${limit}`, "GET", undefined, token);

// Couch API

export const getCouchPlans = (id: string, token: string, limit: number = 30) =>
  request<Cinema[]>(
    `/couches/${id}/plans?limit=${limit}`,
    "GET",
    undefined,
    token,
  );

export const createCouch = (name: string, users: string[], token: string) =>
  request<Couch>("/couches", "POST", { name, users }, token);

export const updateCouch = (
  id: string,
  data: Partial<{ name: string; users: string[] }>,
  token: string,
) => request<Couch>(`/couches/${id}`, "PUT", data, token);

export const getCouch = (id: string, token: string) =>
  request<Couch>(`/couches/${id}`, "GET", undefined, token);

export const getAllCouches = (token: string) =>
  request<Couch[]>("/couches", "GET", undefined, token);
export const checkAdmin = (token: string) =>
  request<void>("/admin", "GET", undefined, token);
