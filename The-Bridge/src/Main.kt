import java.util.*

val commands: List<String> = listOf("SPEED", "SLOW", "WAIT", "JUMP", "UP", "DOWN")

class Road(input: Scanner) {
    var numberBikes = input.nextInt()
    val minimumBikes = input.nextInt()
    val lanes: List<String> = List(4) {input.next()}
    val bikes: MutableList<Bike> = mutableListOf()
    val commandsMade: MutableList<String> = mutableListOf()

    fun isEndOfRoad(): Boolean {
        val activatedBikes = bikes.filter {it.isActivated}
        if (activatedBikes.isNotEmpty()) return activatedBikes[0].y >= lanes[0].length - 1
        return false
    }

    fun isChangeLanePossible(command: String): Boolean {
        val checkLane = if (command == "DOWN") 3 else 0
        return bikes.filter {it.x == checkLane && it.isActivated}.isEmpty()
    }

    fun clearState() {
        bikes.clear()
        commandsMade.clear()
    }
}

class Bike(var speed: Int, input: Scanner) {
    var y = input.nextInt()
    var x = input.nextInt()
    var isActivated = input.nextInt() == 1
    val bikeState = BikeState(speed, x, y, isActivated)

    fun sync(command: String, road: Road) {
        if (!isActivated) return
        isActivated = when(command) {
            "SPEED" -> move(1, false, 0, road)
            "WAIT" -> move(0, false, 0, road)
            "SLOW" -> move(- 1, false, 0, road)
            "JUMP" -> move(0, true, 0, road)
            "UP" -> move(0, false, - 1, road)
            "DOWN" -> move(0, false, 1, road)
            else -> false
        }
        updateState()
    }

    fun move(addSpeed: Int, isJump: Boolean, switch: Int, road: Road): Boolean {
        speed += addSpeed
        if (speed < 0) speed = 0
        val firstY = y + 1
        val lane = road.lanes[x]
        y = listOf(y + speed, lane.length - 1).min()
        if (isJump) return lane[y] != '0'
        else if (switch == 0) return '0' !in lane.slice(firstY..y)
        else {
            if ('0' !in lane.slice(firstY until y) && 0 <= x + switch && x + switch < 4 && '0' !in road.lanes[x + switch].slice(firstY..y)) {
                x += switch
                return true
            }
        }
        return false
    }

    fun updateState() {
        bikeState.speed = speed
        bikeState.x = x
        bikeState.y = y
        bikeState.isActivated = isActivated
    }

    fun restoreState(state: BikeState) {
        speed = state.speed
        x = state.x
        y = state.y
        isActivated = state.isActivated
        updateState()
    }
}

data class BikeState(var speed: Int, var x: Int, var y: Int, var isActivated: Boolean) {
    fun copyState() = BikeState(speed, x, y, isActivated)
}

fun dps(road: Road, extraBike: Int = 0): Boolean {
    if (road.isEndOfRoad()) {
        if (road.commandsMade.isEmpty()) road.commandsMade += "SPEED"
        return true
    }
    for (command in commands) {
        if (command in listOf("DOWN", "UP") && !road.isChangeLanePossible(command)) continue
        val bikesStates = road.bikes.map {it.bikeState.copyState()}
        road.bikes.forEach {it.sync(command, road)}
        road.commandsMade += command
        road.numberBikes = road.bikes.filter {it.isActivated}.size
        if (road.numberBikes >= road.minimumBikes + extraBike && road.bikes.filter {it.isActivated && it.speed == 0}.isEmpty()) {
            if (dps(road, extraBike)) return true
        }
        road.commandsMade.removeLast()
        for (i in 0 until road.bikes.size) {
            road.bikes[i].restoreState(bikesStates[i])
        }
    }
    return false
}

fun main() {
    val input = Scanner(System.`in`)
    val road = Road(input)
    val numberBikes = road.numberBikes

    // game loop
    while (true) {
        road.clearState()
        val bikeSpeed = input.nextInt()
        repeat(numberBikes) {road.bikes += Bike(bikeSpeed, input)}
        if (road.minimumBikes + 1 <= numberBikes && dps(road, 1)) println(road.commandsMade[0])
        else {
            dps(road)
            println(road.commandsMade[0])
        }
    }
}