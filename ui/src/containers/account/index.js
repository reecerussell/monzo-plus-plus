import React, { useState, useEffect } from "react";
import Fetch from "../../utils/fetch";
import * as User from "../../utils/user";
import Index from "../../components/account/index";

const IndexContainer = () => {
	const [userData, setUserData] = useState(null);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState(null);

	const handleFetchUser = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/" + User.GetId(),
			null,
			async (res) => setUserData(await res.json()),
			setError
		);

		setLoading(false);
	};

	useEffect(() => {
		handleFetchUser();
	}, []);

	return <Index data={userData} error={error} loading={loading} />;
};

export default IndexContainer;
