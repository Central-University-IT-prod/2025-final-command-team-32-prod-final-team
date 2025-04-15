// AdminDel.tsx
import * as React from "react";
import { useForm } from "oxyform";
import { Cinema, deleteFilm, searchFilm } from "@/client.js";
import { toast } from "sonner";
import { useAuth } from "@/Auth.js";
import { useNavigate } from "react-router-dom";
import MultipleSelector from "@/components/ui/common/multiple-selector";
import { Button } from "@/components/ui/button";

const AdminDelFilm = () => {
  const navigate = useNavigate();
  const [, , , token] = useAuth();
  const [filmQuery, setFilmQuery] = React.useState("");
  const [foundFilms, setFoundFilms] = React.useState<Cinema[]>([]);
  const [selectedFilm, setSelectedFilm] = React.useState<Cinema | null>(null);

  const form = useForm({
    initialValues: {
      film: null as Cinema | null,
    },
    validation: {
      film: [(value) => !!value || "Please select a film to delete"],
    },
  });
  React.useEffect(() => {
    if (filmQuery.length > 1) {
      searchFilm(filmQuery)
        .then(setFoundFilms)
        .catch((error) => {
          console.error(error);
          toast.error("Error searching films");
        });
    } else {
      setFoundFilms([]);
    }
  }, [filmQuery]);

  const handleDelete = async () => {
    if (!selectedFilm) return;
    console.log(selectedFilm.id);

    try {
      await deleteFilm(selectedFilm.id, token);
      toast.success("Film deleted successfully");
      navigate("/adminPanel");
    } catch (error) {
      console.error(error);
      toast.error("Error deleting film");
    }
  };

  return (
    <div
      className="mx-auto max-w-xl space-y-4 p-4"
      style={{ minWidth: "45vw" }}
    >
      <h1 className="mb-4 text-3xl">Delete Film</h1>

      <div className="space-y-2">
        <label className="block">
          <p className="mb-2 text-xl">Search Film</p>
          <MultipleSelector
            options={foundFilms.map((film) => ({
              value: film.id,
              label: film.name,
            }))}
            value={
              selectedFilm
                ? [
                    {
                      value: selectedFilm.id,
                      label: selectedFilm.name,
                    },
                  ]
                : []
            }
            onValueChange={(options) => {
              const film = foundFilms.find((f) => f.id === options[0]?.value);
              setSelectedFilm(film || null);
            }}
            placeholder="Type to search films..."
            inputValue={filmQuery}
            onInputChange={setFilmQuery}
          />
          {form.errors.film && (
            <div className="error-text">{form.errors.film}</div>
          )}
        </label>

        <Button
          type="button"
          onClick={handleDelete}
          className="bg-destructive hover:bg-destructive/90 mt-10 w-full text-lg"
          disabled={!selectedFilm}
        >
          Delete Selected Film
        </Button>
      </div>
    </div>
  );
};

export default AdminDelFilm;
