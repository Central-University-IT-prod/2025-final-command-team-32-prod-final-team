import React from "react";
import { Button } from "@/components/ui/lib/button.tsx";
import { Input } from "@/components/ui/lib/input.tsx";
import { ErrorWrapper } from "@/ErrorWrapper.jsx";
import { useForm } from "oxyform";
import { addAdminFilm, getGenres } from "@/client.js";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { useAuth } from "@/Auth.js";
import MultipleSelector from "@/components/ui/common/multiple-selector.tsx";

const AdminCreate = () => {
  const navigate = useNavigate();
  const [, , , token] = useAuth();
  const [genreOptions, setGenreOptions] = React.useState([]);
  const [genreQuery, setGenreQuery] = React.useState("");

  // Загрузка жанров с сервера
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

  const form = useForm({
    initialValues: {
      name: "",
      description: "",
      poster_url: "",
      genres: [],
    },
    validation: {
      name: ["required", /^[a-zA-Zа-яА-Я0-9.,!?:;()\s"'«»—–-]+$/],
      description: ["required", /^[a-zA-Zа-яА-Я0-9.,!?:;()\s"'«»—–-]*$/],
      poster_url: [],
      genres: [
        (value) =>
          (Array.isArray(value) && value.length > 0) ||
          "Select at least one genre",
      ],
    },
  });

  // Фильтрация жанров по запросу
  const filteredGenres = genreOptions.filter((genre) =>
    genre.label.toLowerCase().includes(genreQuery.toLowerCase()),
  );

  // Обработчик изменения выбранных жанров
  const handleGenreChange = (newOptions) => {
    form.setValue("genres", newOptions);
  };

  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const handleFormSubmit = async (values) => {
    try {
      setIsSubmitting(true);
      const genresArray = values.genres.map((opt) => opt.value);
      await addAdminFilm(
        values.name,
        values.description,
        values.poster_url,
        genresArray,
        token,
      );
      toast.success("Film added successfully!");
      navigate("/");
    } catch (error) {
      console.error(error);
      toast.error("Error adding film");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form
      style={{ minWidth: "45vw" }}
      className="mx-auto max-w-xl space-y-4 p-4"
      onSubmit={(e) => {
        e.preventDefault();
        form.submit(handleFormSubmit);
      }}
    >
      <div className="space-y-2">
        <h1 className="mb-4 text-3xl">Add your film</h1>
        <label className="block">
          <p className="mb-2 text-xl">Film title</p>
          <ErrorWrapper {...form.register("name")}>
            <Input
              placeholder="Enter your film title"
              disabled={isSubmitting}
            />
          </ErrorWrapper>
        </label>

        <label className="block">
          <p className="mb-2 text-xl">Description</p>
          <ErrorWrapper {...form.register("description")}>
            <Input placeholder="Describe your film" disabled={isSubmitting} />
          </ErrorWrapper>
        </label>

        {/*<label className="block">
          <p className="mb-2 text-xl">Poster URL</p>
          <Input
            OnChange={(x) => console.log(x)}
            type="file"
            placeholder="Add image"
            accept="image/*"
            disabled={isSubmitting}
          />
        </label>*/}

        <label className="block">
          <p className="mb-2 text-xl">Genres</p>
          <MultipleSelector
            options={filteredGenres}
            value={form.values.genres}
            onValueChange={handleGenreChange}
            placeholder="Select genres..."
            disabled={isSubmitting}
            inputValue={genreQuery}
            onInputChange={setGenreQuery}
            hidePlaceholderWhenSelected
          />
          {form.errors.genres && (
            <div className="error-text">{form.errors.genres}</div>
          )}
        </label>
      </div>

      <Button type="submit" className="w-full text-lg" disabled={isSubmitting}>
        {isSubmitting ? "Submitting..." : "Submit"}
      </Button>
    </form>
  );
};

export default AdminCreate;
