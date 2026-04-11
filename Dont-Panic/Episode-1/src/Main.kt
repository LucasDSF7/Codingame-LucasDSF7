import java.util.*

/*
This solution could be easer than it's been done here, but I
think that creating a full simulation is interesting.
 */

class Area(input: Scanner) {
    val nFloors = input.nextInt() // i
    val width = input.nextInt() // j
    val nRounds = input.nextInt()
    val exitFloor = input.nextInt()
    val exitPos = input.nextInt()
    val nTotalClones = input.nextInt()
    val nAdditionalElevators = input.nextInt()
    val nElevators = input.nextInt()
    val elevatorsMap: MutableMap<Int, Int> = mappingElevators(input) // key -> floor to value -> position
    val generatorFloor: Int = input.nextInt()
    val generatorPos: Int = input.nextInt()

    override fun toString(): String {
        return "Number of floors $nFloors\nWidth: $width\nNumber of rounds: $nRounds\n" +
                "Exit Floor: $exitFloor\nExit pos: $exitPos\n" +
                "Number total clones: $nTotalClones\nNumber additional elevators: $nAdditionalElevators\n" +
                "Number Elevators: $nElevators\nElevators map: $elevatorsMap\n" +
                "Generator Floor: $generatorFloor\nGenerator Position: $generatorPos"
    }

    fun mappingElevators(input: Scanner): MutableMap<Int, Int> {
        val auxMap: MutableMap<Int, Int> = mutableMapOf()
        repeat(nElevators) {auxMap[input.nextInt()] = input.nextInt()}
        return auxMap
    }
}

class ImprobabilityDrive(val area: Area) {
    val board = generateBoard()
    var clones: MutableList<Clone> = mutableListOf()
    var leadingClone: Clone? = null

    override fun toString(): String {
        return "${board.reversed().joinToString("\n")}\nLeading $leadingClone\nClones:\n${clones.joinToString("\n")}\n"
    }

    fun generateBoard(): MutableList<MutableList<Char>> {
        val board: MutableList<MutableList<Char>> = mutableListOf()
        for (i in 0 until area.nFloors) {
            val auxList: MutableList<Char> = mutableListOf()
            for (j in 0 until area.width) {
                val char =
                    if (i == area.generatorFloor && j == area.generatorPos) 'G'
                    else if (i in area.elevatorsMap && j == area.elevatorsMap[i]) '^'
                    else if (i == area.exitFloor && j == area.exitPos) 'X'
                    else if (j == 0 || j == area.width - 1) '|'
                    else '_'
                auxList += char
            }
            board += auxList
        }
        return board
    }
}

/*
 Clones don't move when appear from generator
 Before a clone moves it validates if the next position is blocked by a clone
 if so, it changes its direction and move. If the actual position is blocked, when a clone
 emerges from the generator or elevator, by a clone
 the moving clone change its direction and move.
 */

class Clone(generatorFloor: Int, generatorPos: Int, generatorDirection: Int = 1) {
    var floor = generatorFloor
    var pos = generatorPos
    var direction = generatorDirection // Clone always begin facing RIGHT

    override fun toString(): String {
        return "Clone at floor: $floor Position: $pos moving $direction"
    }

    fun move(board: MutableList<MutableList<Char>>) : Boolean {
        val char = board[floor][pos]
        val nextChar = board[floor][pos + direction]
        if (char == '^') {
            floor += 1
            return true
        }
        else if (char == 'B' || nextChar == 'B') {
            direction *= - 1
            pos += direction
            return true
        }
        else if (nextChar == '|') {
            return false
        }
        else pos += direction
        return true
    }

    fun leaderMove(command: String, impDrive: ImprobabilityDrive): Boolean {
        val board = impDrive.board
        val optRange = if (direction == 1) pos until impDrive.area.width else pos downTo 1
        val slice = board[floor].slice(optRange)
        if (command == "BLOCK" && 'X' !in slice && this.getChar(board) != '^') {
            board[floor][pos] = 'B'
            return true
        }
        else {
            // System.err.println(slice)
            if ('^' !in slice && 'X' !in slice && 'B' !in slice && slice.isNotEmpty()) return false
        }
        return move(board)
    }

    fun saveState() = CloneState(floor, pos, direction)

    fun getChar(board: MutableList<MutableList<Char>>) = board[floor][pos]
}

data class CloneState(val floor: Int, val pos: Int, val direction: Int) {
    fun createClone() = Clone(floor, pos, direction)
}

class Simulation(val impDrive: ImprobabilityDrive) {
    val area = impDrive.area
    var turnNumber = 0
    val commands = listOf("WAIT", "BLOCK")
    val commandsMade = mutableListOf<String>()

    fun run(): Boolean {

        if (impDrive.leadingClone?.getChar(impDrive.board) == 'X') return true
        if (turnNumber == area.nRounds) return false
        var blockingLeader: Clone? = null
        if (turnNumber % 3 == 0) {
            val newClone = Clone(area.generatorFloor, area.generatorPos)
            if (impDrive.leadingClone == null) {
                impDrive.leadingClone = newClone
            }
            else impDrive.clones += newClone
        }
        System.err.println("Turn: $turnNumber")
        System.err.println("Commands: $commandsMade")
        System.err.println(impDrive)
        turnNumber += 1
        val leaderState: CloneState? = impDrive.leadingClone?.saveState()
        val cloneStates = impDrive.clones.map {it.saveState()}
        impDrive.clones.forEach { it.move(impDrive.board) }
        for (command in commands) {
            commandsMade += command
            if (impDrive.leadingClone == null) {
                if (run()) return true
            }
            else if (impDrive.leadingClone!!.leaderMove(command, impDrive)) {
                if (command == "BLOCK") {
                    blockingLeader = impDrive.leadingClone
                    impDrive.leadingClone = if (impDrive.clones.isNotEmpty()) impDrive.clones.removeFirst() else null
                }
                if (run()) return true
            }
            if (blockingLeader?.getChar(impDrive.board) == 'B') {
                impDrive.board[blockingLeader.floor][blockingLeader.pos] = '_'
            }
            impDrive.leadingClone = leaderState?.createClone()
            commandsMade.removeLast()
        }
        impDrive.clones = cloneStates.map { it.createClone() }.toMutableList()
        System.err.println("Returning!")
        turnNumber -= 1
        return true
    }
}

fun main() {
    val input = Scanner(System.`in`)
    val area = Area(input)
    val impDrive = ImprobabilityDrive(area)
    System.err.println(area)
    val simulation = Simulation(impDrive)
    simulation.run()
    for (cmd in simulation.commandsMade + "WAIT") {
        println(cmd)
    }
}