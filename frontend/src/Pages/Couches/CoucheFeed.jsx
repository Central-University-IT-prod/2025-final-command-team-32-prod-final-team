import React, { useEffect, useState } from "react";
import { useAuth } from "@/Auth.js";
import { Navigate, useParams } from "react-router";
import { getCouch } from "@/client.js";
import { Explore } from "@/Pages/Explore/Explore.jsx";

export function CoucheFeed() {
  const [info, setInfo] = useState({});

  const params = useParams();
  const id = params.id;

  const [auth, , creds, token] = useAuth(true);

  useEffect(() => {
    if (auth) {
      getCouch(id, token).then((x) => setInfo(x));
    }
  }, []);

  if (!auth) {
    return <Navigate to="/auth" />;
  }

  return <Explore group={id} />;
}
