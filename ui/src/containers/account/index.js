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

		const res = await Fetch(
			"http://localhost:9789/auth/users/" + User.GetId()
		);
		if (res.ok) {
			const data = await res.json();

			if (res.status === 200) {
				setUserData(data);
			} else {
				setError(data.error);
			}
		} else {
			setError(res.statusText);
		}

		setLoading(false);
	};

	useEffect(() => {
		handleFetchUser();
	}, []);

	return <Index data={userData} error={error} loading={loading} />;
};

export default IndexContainer;
