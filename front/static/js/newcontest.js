/**
 * Saves the contest data to localStorage
 * @param {Object} contestData - The contest data to be saved
 */
function saveContestToLocalStorage(contestData) {
    localStorage.setItem('newcontest', JSON.stringify(contestData));
}

/**
 * Retrieves the contest data from localStorage
 * @returns {Object} The contest data, or an empty object if not found
 */
function getContestFromLocalStorage() {
    const contestData = localStorage.getItem('newcontest');
    return contestData ? JSON.parse(contestData) : {};
}

/**
 * Updates the contest data in localStorage
 * @param {Object} updatedContestData - The updated contest data
 */
function updateContestInLocalStorage(updatedContestData) {
    saveContestToLocalStorage(updatedContestData);
}

/**
 * Deletes the contest data from localStorage
 */
function deleteContestFromLocalStorage() {
    localStorage.removeItem('newcontest');
}

/**
 * Adds a new outcome to the contest data in localStorage
 * @param {Object} newOutcome - The new outcome to be added
 */
function addOutcomeToLocalStorage(newOutcome) {
    const contestData = getContestFromLocalStorage();
    contestData.outcomes.push(newOutcome);
    saveContestToLocalStorage(contestData);
}

/**
 * Adds a new game to the contest data in localStorage
 * @param {Object} newGame - The new game to be added
 */
function addGameToLocalStorage(newGame) {
    const contestData = getContestFromLocalStorage();
    contestData.games.push(newGame);
    saveContestToLocalStorage(contestData);
}

/**
 * deltes all game to the contest data in localStorage
 */
function deleteGamesFromLocalStorage() {
    const contestData = getContestFromLocalStorage();
    contestData.games = [];
    saveContestToLocalStorage(contestData);
}

function getISOOffsetString(date) {

    const offsetMinutes = date.getTimezoneOffset();
    const offsetHours = Math.floor(Math.abs(offsetMinutes) / 60);
    const offsetMinutesRemaining = Math.abs(offsetMinutes) % 60;
    const offsetSign = offsetMinutes >= 0 ? '+' : '-';

    return `${offsetSign}${String(offsetHours).padStart(2, '0')}:${String(offsetMinutesRemaining).padStart(2, '0')}`;
}
