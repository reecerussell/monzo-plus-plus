import React, { useState } from "react";
import * as User from "../../utils/user";
import Fetch from "../../utils/fetch";
import Delete from "../../components/account/delete";

const DeleteContainer = () => {
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);
	const [redirect, setRedirect] = useState(null);

	const handleDelete = async e => {
		e.preventDefault();

		if (loading) {
			return;
		}

		setLoading(true);

		try {
			const res = await Fetch(
				"http://localhost:9789/auth/users/" + User.GetId(),
				{
					method: "DELETE",
				}
			);

			if (res.status == 200) {
				setError(null);
				User.Logout();
				setRedirect("/");
			} else {
				const data = await res.json();
				setError(data.error);
			}
		} catch {
			setError(
				"It seems like you don't have connection to the internet. Try again later!"
			);
		}

		setLoading(false);
	};

	return (
		<Delete
			handleSubmit={handleDelete}
			error={error}
			loading={loading}
			redirect={redirect}
		/>
	);
};

export default DeleteContainer;
