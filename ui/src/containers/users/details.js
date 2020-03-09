import React, { useState, useEffect } from "react";
import Fetch from "../../utils/fetch";
import Details from "../../components/users/details";

const DetailsContainer = ({ id }) => {
	const [details, setDetails] = useState();
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);

	const fetchUserDetails = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		const res = await Fetch("http://localhost:9789/auth/users/" + id);

		if (res.ok) {
			const data = await res.json();

			if (res.status === 200) {
				setDetails(data);
			} else {
				setError(data.error);
			}
		} else {
			setError(res.statusText);
		}

		setLoading(false);
	};

	useEffect(() => {
		fetchUserDetails();
	}, []);

	return <Details loading={loading} error={error} {...details} />;
};

export default DetailsContainer;
