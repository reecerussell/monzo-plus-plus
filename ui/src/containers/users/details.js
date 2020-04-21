import React, { useState, useEffect } from "react";
import Fetch from "../../utils/fetch";
import Details from "../../components/users/details";

const DetailsContainer = ({ id }) => {
	const [details, setDetails] = useState();
	const [readonly, setReadonly] = useState(true);
	const [error, setError] = useState(null);
	const [success, setSuccess] = useState(null);
	const [loading, setLoading] = useState(false);

	const fetchUserDetails = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/" + id,
			null,
			async (res) => setDetails(await res.json()),
			setError
		);

		setLoading(false);
	};

	const handleSubmit = async (e) => {
		e.preventDefault();

		if (loading || readonly) {
			return;
		}

		setLoading(true);

		await Fetch(
			"auth/users",
			{
				method: "PUT",
				body: JSON.stringify(details),
			},
			async () => {
				setReadonly(true);
				setSuccess("Updated successfully!");
			},
			setError
		);

		setLoading(false);
	};

	const handleUpdate = (e) => {
		const { name, value } = e.target;
		const data = { ...details };
		data[name] = value;
		setDetails(data);
	};

	const toggleMode = (e) => {
		if (e) e.preventDefault();

		setReadonly(!readonly);
	};

	useEffect(() => {
		fetchUserDetails();
	}, []);

	return (
		<Details
			loading={loading}
			error={error}
			handleSubmit={handleSubmit}
			handleUpdate={handleUpdate}
			toggleMode={toggleMode}
			data={details}
			readonly={readonly}
			success={success}
			id={id}
		/>
	);
};

export default DetailsContainer;
