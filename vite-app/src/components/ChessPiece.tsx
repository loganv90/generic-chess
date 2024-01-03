import React from 'react'

const divStyle: React.CSSProperties = {
    width: '100%',
    height: '100%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
}

const imgStyle: React.CSSProperties = {
    imageRendering: 'pixelated',
    objectFit: 'contain',
    width: '80%',
    height: '80%',
}

const getImgSrc = (type: string, color: string): string => {
    if (color === 'white') {
        switch (type) {
            case 'P':
                return '/src/assets/wp.png'
            case 'R':
                return '/src/assets/wr.png'
            case 'N':
                return '/src/assets/wn.png'
            case 'B':
                return '/src/assets/wb.png'
            case 'Q':
                return '/src/assets/wq.png'
            case 'K':
                return '/src/assets/wk.png'
            default:
                return ''
        }
    } else if (color === 'black') {
        switch (type) {
            case 'P':
                return '/src/assets/bp.png'
            case 'R':
                return '/src/assets/br.png'
            case 'N':
                return '/src/assets/bn.png'
            case 'B':
                return '/src/assets/bb.png'
            case 'Q':
                return '/src/assets/bq.png'
            case 'K':
                return '/src/assets/bk.png'
            default:
                return ''
        }
    } else {
        return ''
    }
}

const ChessPiece = ({
    type,
    color,
}: {
    type: string,
    color: string,
}): JSX.Element => {
    const imgSrc = getImgSrc(type, color)
    return (
        <div style={divStyle} >
            {imgSrc ? <img src={imgSrc} style={imgStyle} /> : type + color}
        </div>
    )
}

export default ChessPiece
