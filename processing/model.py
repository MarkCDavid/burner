from dataclasses import dataclass
from typing import Dict


@dataclass
class LeafBlock:
    id: int
    depth: int

    @classmethod
    def from_values(cls, values):
        return cls(id=int(values[0]), depth=int(values[1]))


@dataclass
class Block:
    id: int
    previousBlockId: int
    minedBy: int
    depth: int
    startedAt: float
    finishedAt: float
    abandoned: bool
    transactions: int
    blockType: int

    def duration(self):
        return self.finishedAt - self.startedAt

    def duration_since_previous(self, blocks: "Dict[int, Block]"):
        if self.previousBlockId == 0:
            return self.startedAt
        return self.startedAt - blocks[self.previousBlockId].finishedAt

    def duration_since_previous_type(self, blocks: "Dict[int, Block]"):
        if self.previousBlockId == 0:
            return self.startedAt

        previous = blocks[self.previousBlockId]
        while previous.previousBlockId != 0:
            if previous.blockType == self.blockType:
                return self.startedAt - previous.finishedAt
            previous = blocks[previous.previousBlockId]

    @classmethod
    def from_values(cls, values):
        return cls(
            id=int(values[0]),
            previousBlockId=int(values[1]),
            minedBy=int(values[2]),
            depth=int(values[3]),
            startedAt=float(values[4]),
            finishedAt=float(values[5]),
            abandoned=bool(values[6]),
            transactions=int(values[7]),
            blockType=int(values[8]),
        )


@dataclass
class BlockMiningAverages:
    blockType: int
    blockCount: int
    averageBlockTime: float
    minBlockTime: float
    maxBlockTime: float

    @classmethod
    def from_values(cls, values):
        return cls(
            blockType=int(values[0]),
            blockCount=int(values[1]),
            averageBlockTime=float(values[2]),
            minBlockTime=float(values[3]),
            maxBlockTime=float(values[4]),
        )


@dataclass
class PPOBPricing:
    id: int
    timestamp: float
    nodeId: int
    price: float

    @classmethod
    def from_values(cls, values):
        return cls(
            id=int(values[0]),
            timestamp=float(values[1]),
            nodeId=int(values[2]),
            price=float(values[3]),
        )
