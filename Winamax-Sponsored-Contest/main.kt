import java.util.Scanner

class Board(input: Scanner) {
    val width = input.nextInt()
    val height = input.nextInt()
    val board = createBoard(input)
    val hazardsPos: MutableList<List<Int>> = mutableListOf()
    val balls: MutableList<Ball> = mutableListOf()

    init {
        findBallsHazard()
    }

    private fun createBoard(input: Scanner): MutableList<MutableList<Char>> {
        val board: MutableList<MutableList<Char>> = mutableListOf()
        for (i in 0 until height) {
            val row = input.next()
            val auxList: MutableList<Char> = mutableListOf()
            for (col in row) {
                auxList += col
            }
            board += auxList
        }
        return board
    }

    private fun findBallsHazard() {
        for (i in 0 until height) {
            for (j in 0 until width) {
                val char = board[i][j]
                if (char == 'X') {
                    hazardsPos += listOf(i, j)
                }
                else if (char.isDigit()) {
                    balls += Ball(i, j, char.digitToInt(), this)
                }
            }
        }
    }
}

class Ball(var i: Int, var j: Int, var shotCount: Int, val board: Board) {
    val initPos = listOf(i, j)
    val moves: Map<Char, List<Int>> = mapOf(
        '>' to listOf(0, 1),
        '<' to listOf(0, - 1),
        'v' to listOf(1, 0),
        '^' to listOf(- 1, 0)
    )

    fun move(move: Char): Boolean {
        val dir = moves[move]!!
        val newI = i + dir[0] * shotCount
        val newJ = j + dir[1] * shotCount
        if (isValidMove(newI, newJ)) {
            if (changeBoard(i, newI, j, newJ, move)) {
                i = newI
                j = newJ
                shotCount -= 1
                return true
            }
        }
        return false
    }

    fun unMove(move: Char) {
        // println("Calling unMove.")
        shotCount += 1
        val dir = moves[move]!!
        val oldI = i - dir[0] * shotCount
        val oldJ = j - dir[1] * shotCount
        unChangeBoard(i, oldI, j, oldJ)
        i = oldI
        j = oldJ
    }

    fun isValidMove(newI: Int, newJ: Int): Boolean {
        return (0 <= newI && newI < board.height) && (0 <= newJ && newJ < board.width) && board.board[newI][newJ] !in listOf('X', 'B') + moves.keys && !board.board[newI][newJ].isDigit()
    }

    fun changeBoard(i: Int, newI: Int, j: Int, newJ: Int, move: Char): Boolean {
        val xRange = getRange(i, newI)
        val yRange = getRange(j, newJ)
        var lastX = i
        var lastY = j
        for (x in xRange) {
            for (y in yRange) {
                val char = board.board[x][y]
                if ((char in moves.keys) || ((x != newI || y != newJ) && char in listOf('H', 'B')) || (char.isDigit() && (x != i || y != j))) {
                    unChangeBoard(i, lastX, j, lastY)
                    return false
                }
                board.board[x][y] = move
                lastX = x
                lastY = y
            }
        }
        return true
    }

    fun unChangeBoard(i: Int, oldI: Int, j: Int, oldJ: Int) {
        val xRange = unGetRange(i, oldI)
        val yRange = unGetRange(j, oldJ)
        for (x in xRange) {
            for (y in yRange) {
                val char = board.board[x][y]
                val newChar = if (listOf(x, y) in board.hazardsPos) 'X' else if (listOf(x, y) == this.initPos) this.shotCount.digitToChar() else if (char == 'H') 'H' else '.'
                if (listOf(x, y) == this.initPos) {
                    println("Called from unChangeBoard. newChar = $newChar")
                }
                board.board[x][y] = newChar
            }
        }
    }

    fun getPos(): Char {
        return board.board[i][j]
    }

    fun getRange(from: Int, to: Int): IntProgression {
        return if (from == to) from..to
        else if (from < to) from until to else from downTo to + 1
    }

    fun unGetRange(from: Int, to: Int): IntProgression {
        return if (from == to) from..to
        else if (from < to) from  until to + 1 else from - 1 downTo to
    }
}

fun dps(board: Board): Boolean {
    for (ball in board.balls) {
        if (ball.getPos() == 'H') {
            board.board[ball.i][ball.j] = 'B'
            board.balls -= ball
            if (dps(board)) return true
            board.balls.add(0, ball)
            board.board[ball.i][ball.j] = 'H'
            return false
        }
        else if (ball.shotCount == 0) return false
        for (move in ball.moves.keys) {
            for (row in board.board) {
                println(row)
            }
            println("$move------- ${ball.initPos}")
            if (ball.move(move)) {
                println("now [${ball.i}, ${ball.j}]")
                if (dps(board)) return true
                ball.unMove(move)
                for (row in board.board) {
                    println(row)
                }
                println("now [${ball.i}, ${ball.j}]\n")
            }
        }
        return false
    }
    if (board.balls.isEmpty()) {
        for (row in board.board) {
            println(row.map { if (it !in listOf('<', '>', '^', 'v')) '.' else it }.joinToString(""))
        }
        return true
    }
    return true
}



fun main() {
    val input = Scanner(System.`in`)
    val board = Board(input)
    dps(board)
}