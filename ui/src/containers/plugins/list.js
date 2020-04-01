import React from "react";
import useListState from "../listState";
import List from "../../components/plugins/list";

const ListContainer = () => {
	const [
		plugins,
		loading,
		error,
		searchTerm,
		updateSearchTerm,
		reload,
	] = useListState("api/plugins/");

	const handleUpdateSearch = e => updateSearchTerm(e.target.value);

	const handleSearch = async e => {
		e.preventDefault();

		await reload();
	};

	return (
		<List
			plugins={plugins}
			loading={loading}
			error={error}
			searchTerm={searchTerm}
			handleUpdateSearch={handleUpdateSearch}
			handleSearch={handleSearch}
		/>
	);
};

export default ListContainer;
