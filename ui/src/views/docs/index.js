import React from "react";
import SwaggerUI from "swagger-ui-react";
import "swagger-ui-react/swagger-ui.css";

import spec from "../../docs.js";

const Index = () => <SwaggerUI spec={spec} />;

export default Index;
