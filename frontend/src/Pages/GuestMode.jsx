import React from "react";
import { signUp } from "@/client.js";
import { toast } from "sonner";
import { Button } from "@/components/ui/lib/button.tsx";
import { useNavigate } from "react-router";
import { useLogin } from "@/Auth.js";

const GuestMode = ({ isSubmitting, setIsSubmitting }) => {
  const navigate = useNavigate();
  const login = useLogin();

  function generateName() {
    const letters = "abcdefghijklmnopqrstuvwxyz";
    let name = "";
    for (let i = 0; i < 4; i++) {
      name += letters[Math.floor(Math.random() * letters.length)];
    }
    return name;
  }

  const guestMode = async () => {
    const guestName = "guest-" + generateName();
    try {
      setIsSubmitting(true);
      const response = await signUp(guestName, "password");
      const tokenUntil = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000);
      login(response.token, tokenUntil, { username: guestName });
      toast.success("Registration successful!");
      navigate("/");
    } finally {
      setIsSubmitting(false);
    }
  };
  return (
    <Button
      type="button"
      disabled={isSubmitting}
      className="mt-2 w-full text-lg"
      onClick={guestMode}
    >
      Guest Mode
    </Button>
  );
};

export default GuestMode;
