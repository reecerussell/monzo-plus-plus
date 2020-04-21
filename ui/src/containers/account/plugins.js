import React, { useState, useEffect } from "react";
import Fetch from "../../utils/fetch";
import Plugins from "../../components/account/plugins";
import * as User from "../../utils/user";

const PluginsContainer = () => {
	const [plugins, setPlugins] = useState([]);
	const [error, setError] = useState(null);
	const [loading, setLoading] = useState(false);
	const [searchTerm, setSearchTerm] = useState("");
	const [showMore, setShowMore] = useState(false);

	const toggleMore = () => setShowMore(!showMore);

	const updateSearchTerm = (e) => setSearchTerm(e.target.value);

	const fetchPlugins = async () => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/plugins/?term=" + searchTerm,
			null,
			async (res) => setPlugins(await res.json()),
			setError
		);

		setLoading(false);
	};

	const handleEnablePlugin = async (pluginId) => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/plugin",
			{
				method: "POST",
				body: JSON.stringify({
					pluginId: pluginId,
					userId: User.GetId(),
				}),
			},
			fetchPlugins,
			setError
		);

		setLoading(false);
	};

	const handleDisablePlugin = async (pluginId) => {
		if (loading) {
			return;
		}

		setLoading(true);

		await Fetch(
			"api/auth/users/plugin",
			{
				method: "DELETE",
				body: JSON.stringify({
					pluginId: pluginId,
					userId: User.GetId(),
				}),
			},
			fetchPlugins,
			setError
		);

		setLoading(false);
	};

	const handleSearch = async (e) => {
		e.preventDefault();
		await fetchPlugins();
	};

	useEffect(() => {
		fetchPlugins();
	}, []);

	return (
		<Plugins
			error={error}
			loading={loading}
			plugins={plugins}
			searchTerm={searchTerm}
			updateSearchTerm={updateSearchTerm}
			showMore={showMore}
			toggleMore={toggleMore}
			handleSearch={handleSearch}
			handleEnablePlugin={handleEnablePlugin}
			handleDisablePlugin={handleDisablePlugin}
		/>
	);
};

export default PluginsContainer;
