import { CardLike } from "@/Pages/Explore/Card.jsx";
import React, { useEffect, useState } from "react";
import { NavLink, useSearchParams } from "react-router";
import { getGenres, searchFilm } from "@/client.js";
import { Input } from "@/components/ui/lib/input.tsx";
import MultipleSelector from "@/components/ui/common/multiple-selector.js";
import { toast } from "sonner";

export function Search() {
  const [data, setData] = useState([]);
  const [debouncedQuery, setDebouncedQuery] = useState("");

  const [searchParams, setSearchParams] = useSearchParams();

  const [query, setQuery] = useState(searchParams.get("query") ?? "");

  const [tags, setTags] = useState(searchParams.get("tags")?.split(";") ?? []);

  useEffect(() => {
    setSearchParams({ query: debouncedQuery, tags: tags?.join(";") });
  }, [debouncedQuery, tags]);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedQuery(query);
    }, 500);

    return () => clearTimeout(handler);
  }, [query]);

  useEffect(() => {
    searchFilm(debouncedQuery, tags).then(setData);
  }, [debouncedQuery, tags]);

  const [genreOptions, setGenreOptions] = React.useState([]);
  const [genreQuery, setGenreQuery] = React.useState("");

  React.useEffect(() => {
    getGenres()
      .then((genres) => {
        const options = genres.map((genre) => ({
          value: genre,
          label: genre,
        }));
        setGenreOptions(options);
      })
      .catch((error) => {
        console.error(error);
        toast.error("Error fetching genres");
      });
  }, []);

  const filteredGenres = genreOptions.filter((genre) =>
    genre.label.toLowerCase().includes(genreQuery.toLowerCase()),
  );

  return (
    <div className="mx-auto w-[min(1500px,100%)] px-[2rem]">
      <h1
        className={"text-extrabold mb-[1rem] text-center text-4xl md:text-left"}
      >
        Search
      </h1>
      <Input
        value={query}
        onInput={(x) => setQuery(x.target.value)}
        className="mx-auto"
        style={{ fontSize: "2rem", height: "2em", width: "max(100%,20rem)" }}
        placeholder="Interstellar"
      />
      <p className="mb-2 text-xl">Genres</p>
      <MultipleSelector
        options={filteredGenres}
        value={tags.map((x) => ({ value: x, label: x }))}
        onValueChange={(x) => setTags(x.map((y) => y.label))}
        placeholder="Select genres..."
        inputValue={genreQuery}
        onInputChange={setGenreQuery}
        hidePlaceholderWhenSelected
      />
      <div
        className="mt-[2rem] grid gap-[2rem]"
        style={{
          gridTemplateColumns: "repeat(auto-fill, minmax(350px, 1fr))",
        }}
      >
        {data.map((item) => (
          <CardLike key={item.id} info={item} />
        ))}
      </div>
      <div className="mt-[2rem] w-full text-center text-[2rem]">
        Nothing else found here :(
        <br />
        <NavLink to="/create" className="button underline">
          You can add your own movie
        </NavLink>
      </div>
    </div>
  );
}
