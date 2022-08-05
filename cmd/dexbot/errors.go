/*
   Copyright 2022 dexenrage

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

const (
	errCreateDir = `Can't create %s directory: %v`
	errZapInit   = `Can't initialize zap logger: %v`

	errLoadingConfig = `error loading config: %v`

	errGetWorkDir = `Error getting workdir: %v`

	errDatabaseConnection = `Database connection error: %v`
	errParsingDatabaseURI = `Error parsing database URI: %v`

	errURLNotAllowed = `[ %s ]: URL not allowed`
	errOutOfRange    = `[ %s ]: index out of range`

	errHTTPStatusCode   = `Parse: HTTP status code for %s is not 200. Status code: %d`
	errCannotParsePrice = `Parse: Cannot parse price for %s`
)
