import { describe, it, expect } from 'vitest'
import { Board } from '/src/utils/chess/Board'
import { Pawn, Knight, Bishop, Rook, Queen, King } from '/src/utils/chess/Piece'

describe('Board', () => {
    it('should have default pieces', () => {
        const board = new Board()
        for (let i=0; i<8; i+=7) {
            expect(board.getPiece(0, i).constructor).toBe(Rook)
            expect(board.getPiece(1, i).constructor).toBe(Knight)
            expect(board.getPiece(2, i).constructor).toBe(Bishop)
            expect(board.getPiece(3, i).constructor).toBe(Queen)
            expect(board.getPiece(4, i).constructor).toBe(King)
            expect(board.getPiece(5, i).constructor).toBe(Bishop)
            expect(board.getPiece(6, i).constructor).toBe(Knight)
            expect(board.getPiece(7, i).constructor).toBe(Rook)
        }
        for (let i=0; i<8; i++) {
            expect(board.getPiece(i, 1).constructor).toBe(Pawn)
            expect(board.getPiece(i, 6).constructor).toBe(Pawn)
            expect(board.getPiece(i, 2)).toBe(null)
            expect(board.getPiece(i, 3)).toBe(null)
            expect(board.getPiece(i, 4)).toBe(null)
            expect(board.getPiece(i, 5)).toBe(null)
            for (let j=0; j<2; j++) {
                expect(board.getPiece(0, j).color).toBe('black')
            }
            for (let j=6; j<8; j++) {
                expect(board.getPiece(0, j).color).toBe('white')
            }
        }
    })

    describe('getMoves', () => {
        it('should return correct moves for unmoved pawn', () => {
            const board = new Board()
    
            const eMoves = [
                {x: 4, y: 5},
                {x: 4, y: 4, options: {enPassant: {x: 4, y: 5}}}
            ]
            const aMoves = board.getMoves(4, 6)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for moved pawn', () => {
            const board = new Board()
            board.getPiece(4, 6).move()
            
            const eMoves = [
                {x: 4, y: 5}
            ]
            const aMoves = board.getMoves(4, 6)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for pawn en passant', () => {
            const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/3pP3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
            board.getPiece(4, 3).move()
            board.getPiece(3, 3).move()
            board.enPassant['black'] = {x: 3, y: 2, xPiece: 3, yPiece: 3}
    
            const eMoves = [
                {x: 4, y: 2},
                {x: 3, y: 2, options: {canEnPassant: true}}
            ]
            const aMoves = board.getMoves(4, 3)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for pawn capture', () => {
            const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/3p4/4P3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
            board.getPiece(4, 3).move()
            board.getPiece(3, 2).move()
            
            const eMoves = [
                {x: 4, y: 2},
                {x: 3, y: 2}
            ]
            const aMoves = board.getMoves(4, 3)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })

        it('should return correct moves for pawn promotion', () => {
            const board = new Board(0, 0, ['white', 'black'], 'rnbqkbn1/pppppppP/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
            board.getPiece(7, 1).move()

            const eMoves = [
                {x: 7, y: 0, options: {canPromote: true}},
                {x: 6, y: 0, options: {canPromote: true}}
            ]
            const aMoves = board.getMoves(7, 1)

            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for unmoved knight', () => {
            const board = new Board()
    
            const eMoves = [
                {x: 0, y: 5},
                {x: 2, y: 5}
            ]
            const aMoves = board.getMoves(1, 7)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for moved knight', () => {
            const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/4N3/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
            const eMoves = [
                {x: 3, y: 1},
                {x: 5, y: 1},
                {x: 2, y: 2},
                {x: 6, y: 2},
                {x: 2, y: 4},
                {x: 6, y: 4},
                {x: 3, y: 5},
                {x: 5, y: 5}
            ]
            const aMoves = board.getMoves(4, 3)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for unmoved bishop', () => {
            const board = new Board()
    
            const eMoves = []
            const aMoves = board.getMoves(2, 7)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for moved bishop', () => {
            const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/6B1/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
            const eMoves = [
                {x: 7, y: 4},
                {x: 7, y: 2},
                {x: 5, y: 2},
                {x: 4, y: 1},
                {x: 5, y: 4},
                {x: 4, y: 5}
            ]
            const aMoves = board.getMoves(6, 3)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for unmoved rook', () => {
            const board = new Board()
    
            const eMoves = []
            const aMoves = board.getMoves(0, 7)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for moved rook', () => {
            const board = new Board(0, 0, ['white','black'], 'rnbqkbnr/pppppppp/8/8/R7/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
            const eMoves = [
                {x: 0, y: 5},
                {x: 0, y: 3},
                {x: 0, y: 2},
                {x: 0, y: 1},
                {x: 1, y: 4},
                {x: 2, y: 4},
                {x: 3, y: 4},
                {x: 4, y: 4},
                {x: 5, y: 4},
                {x: 6, y: 4},
                {x: 7, y: 4}
            ]
            const aMoves = board.getMoves(0, 4)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for unmoved queen', () => {
            const board = new Board()
    
            const eMoves = []
            const aMoves = board.getMoves(3, 7)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for moved queen', () => {
            const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/8/3Q4/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
            const eMoves = [
                {x: 4, y: 5},
                {x: 4, y: 3},
                {x: 4, y: 4},
                {x: 5, y: 4},
                {x: 5, y: 2},
                {x: 6, y: 1},
                {x: 6, y: 4},
                {x: 7, y: 4},
                {x: 2, y: 5},
                {x: 2, y: 3},
                {x: 2, y: 4},
                {x: 1, y: 4},
                {x: 1, y: 2},
                {x: 0, y: 1},
                {x: 0, y: 4},
                {x: 3, y: 5},
                {x: 3, y: 3},
                {x: 3, y: 2},
                {x: 3, y: 1}
            ]
            const aMoves = board.getMoves(3, 4)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for unmoved king', () => {
            const board = new Board()
    
            const eMoves = []
            const aMoves = board.getMoves(4, 7)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    
        it('should return correct moves for moved king', () => {
            const board = new Board(0, 0, ['white', 'black'], 'rnbqkbnr/pppppppp/8/8/4K3/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1')
    
            const eMoves = [
                {x: 5, y: 5},
                {x: 5, y: 4},
                {x: 5, y: 3},
                {x: 3, y: 4},
                {x: 3, y: 5},
                {x: 4, y: 5},
                {x: 3, y: 3},
                {x: 4, y: 3}
            ]
            const aMoves = board.getMoves(4, 4)
    
            expect(aMoves).toEqual(expect.arrayContaining(eMoves))
            expect(aMoves.length).toBe(eMoves.length)
        })
    })

    // describe('movePiece', () => {
    //     it('should move piece when player turn', () => {
    //         const board = new Board()

    //         const fromX = 4
    //         const fromY = 6
    //         const toX = 4
    //         const toY = 4
    //         const piece = board.squares[fromY][fromX].piece
    //         const res = board.movePiece(fromX, fromY, toX, toY)

    //         expect(res).toBe(true)
    //         expect(board.squares[fromY][fromX].piece).toBe(null)
    //         expect(board.squares[toY][toX].piece).toBe(piece)
    //         expect(board.currentPlayer).toBe(1)
    //     })
    
    //     it('should not move piece when not player turn', () => {
    //         const board = new Board()

    //         const fromX = 4
    //         const fromY = 1
    //         const toX = 4
    //         const toY = 3
    //         const piece = board.squares[fromY][fromX].piece
    //         const res = board.movePiece(fromX, fromY, toX, toY)

    //         expect(res).toBe(false)
    //         expect(board.squares[fromY][fromX].piece).toBe(piece)
    //         expect(board.squares[toY][toX].piece).toBe(null)
    //         expect(board.currentPlayer).toBe(0)
    //     })
    
    //     it('should not move piece when invalid from', () => {
    //         const board = new Board()

    //         const fromX = 4
    //         const fromY = 5
    //         const toX = 4
    //         const toY = 4
    //         const piece = board.squares[fromY][fromX].piece
    //         const res = board.movePiece(fromX, fromY, toX, toY)

    //         expect(res).toBe(false)
    //         expect(board.squares[fromY][fromX].piece).toBe(piece)
    //         expect(board.squares[toY][toX].piece).toBe(null)
    //         expect(board.currentPlayer).toBe(0)
    //     })
    
    //     it('should not move piece when invalid to', () => {
    //         const board = new Board()

    //         const fromX = 4
    //         const fromY = 6
    //         const toX = 4
    //         const toY = 3
    //         const piece = board.squares[fromY][fromX].piece
    //         const res = board.movePiece(fromX, fromY, toX, toY)

    //         expect(res).toBe(false)
    //         expect(board.squares[fromY][fromX].piece).toBe(piece)
    //         expect(board.squares[toY][toX].piece).toBe(null)
    //         expect(board.currentPlayer).toBe(0)
    //     })

    //     it('should record en passant and clear it after next move', () => {
    //         const board = new Board()

    //         const fromX = 4
    //         const fromY = 6
    //         const toX = 4
    //         const toY = 4
    //         const fromXb = 4
    //         const fromYb = 1
    //         const toXb = 4
    //         const toYb = 3

    //         const piece = board.squares[fromY][fromX].piece
    //         board.movePiece(fromX, fromY, toX, toY)
    //         expect(piece.enPassant).toEqual({x: 4, y: 5})
    //         board.movePiece(fromXb, fromYb, toXb, toYb)
    //         expect(piece.enPassant).toEqual(null)
    //     })
    // })
})
