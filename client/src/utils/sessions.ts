export async function getSession() {
  let jsonData: { name: string } = { name: "" };
  try {
    const response = await fetch("http://localhost:8080/get-session", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      console.error("Could not get session data");
      return { name: "" };
    }

    const data = (await response.json()) as {
      name: string;
    };

    jsonData = data;
  } catch (e) {
    console.error(e);
  }

  return jsonData;
}

export async function doesSessionExist(): Promise<boolean> {
  try {
    const res = await fetch("http://localhost:8080/does-session-exist", {
      method: "POST",
    });

    if (!res.ok) {
      return false;
    }
  } catch (e) {
    console.error(e);
  }

  return true;
}
